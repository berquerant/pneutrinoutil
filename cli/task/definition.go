package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"gopkg.in/yaml.v3"
)

func NewGenerator(dir *Dir, c *ctl.Config, play, hook string) *Generator {
	return &Generator{
		dir:  dir,
		c:    c,
		play: play,
		hook: hook,
	}
}

type Generator struct {
	dir  *Dir
	c    *ctl.Config
	play string
	hook string
}

func (g Generator) dyldLibraryPath() string {
	x, _ := filepath.Abs(g.dir.NeutrinoDir())
	xs := []string{}
	for _, a := range []string{
		os.Getenv("DYLD_LIBRARY_PATH"),
		x + "/bin",
	} {
		if a != "" {
			xs = append(xs, a)
		}
	}
	return strings.Join(xs, ":") + ":"
}

func (g Generator) env() execx.Env {
	e := g.c.Env()
	e.Set("ResultDestDir", g.dir.ResultDestDir())
	e.Set("Score", g.c.Score)
	e.Set("Play", g.play)
	e.Set("Hook", g.hook)
	e.Set("DYLD_LIBRARY_PATH", g.dyldLibraryPath())
	e.Set("HOME", os.Getenv("HOME"))
	e.Merge(g.dir.Env())
	return e
}

func (g Generator) ExecutableTasks() *execx.ExecutableTasks {
	tasks := execx.NewTasks().
		Add(execx.NewTask(
			"init",
			fmt.Sprintf(`export DYLD_LIBRARY_PATH="$DYLD_LIBRARY_PATH"
chmod 755 %[1]s/*
xattr -dr com.apple.quarantine "%[1]s"
cp -f "${Score}" "%[2]s/"
mkdir -p "${ResultDestDir}"`,
				g.dir.BinDir(),
				g.dir.MusicXMLDir(),
			))).
		Add(execx.NewTask(
			"MusicXMLtoLabel",
			fmt.Sprintf(
				`%s/musicXMLtoLabel \
  "%s/${BASENAME}.musicxml" \
  "%s/${BASENAME}.lab" \
  "%s/${BASENAME}.lab"`,
				g.dir.BinDir(),
				g.dir.MusicXMLDir(),
				g.dir.FullDir(),
				g.dir.MonoDir(),
			))).
		Add(execx.NewTask(
			"NEUTRINO",
			fmt.Sprintf(
				`%[1]s/NEUTRINO \
  "%[2]s/${BASENAME}.lab" \
  "%[3]s/${BASENAME}.lab" \
  "%[4]s/${BASENAME}.f0" \
  "%[4]s/${BASENAME}.melspec" \
  "%[5]s/${ModelDir}/" \
  -w "%[4]s/${BASENAME}.mgc" \
  "%[4]s/${BASENAME}.bap" \
  -n ${NumParallel} \
  -o ${NumThreads} \
  -k ${StyleShift} \
  -d ${InferenceMode} \
  -r ${RandomSeed} \
  -i "%[4]s/${BASENAME}.trace" \
  -t`,
				g.dir.BinDir(),
				g.dir.FullDir(),
				g.dir.TimingDir(),
				g.dir.OutputDir(),
				g.dir.ModelDir(),
			))).
		Add(execx.NewTask(
			"NSF",
			fmt.Sprintf(
				`%[1]s/NSF \
  "%[2]s/${BASENAME}.f0" \
  "%[2]s/${BASENAME}.melspec" \
  "%[3]s/${ModelDir}/${NsfModel}.bin" \
  "%[2]s/${BASENAME}.wav" \
  -l "%[4]s/${BASENAME}.lab" \
  -n ${NumParallel} \
  -p ${NumThreads} \
  -s ${SamplingFreq} \
  -f ${PitchShiftNsf} \
  -t`,
				g.dir.BinDir(),
				g.dir.OutputDir(),
				g.dir.ModelDir(),
				g.dir.TimingDir(),
			))).
		Add(execx.NewTask(
			"WORLD",
			fmt.Sprintf(
				`%[1]s/WORLD \
  "%[2]s/${BASENAME}.f0" \
  "%[2]s/${BASENAME}.mgc" \
  "%[2]s/${BASENAME}.bap" \
  "%[2]s/${BASENAME}_world.wav" \
  -f ${PitchShiftWorld} \
  -m ${FormantShift} \
  -p ${SmoothPicth} \
  -c ${SmoothFormant} \
  -b ${EnhanceBreathiness} \
  -n ${NumThreads} \
  -t`,
				g.dir.BinDir(),
				g.dir.OutputDir(),
			))).
		Add(execx.NewTask(
			"cleanup",
			fmt.Sprintf(
				`cp %[1]s/${BASENAME}.* "%[2]s/${BASENAME}.musicxml" "%[1]s/${BASENAME}_world.wav" "${ResultDestDir}/"
cat <<EOS > "${ResultDestDir}/config.yml"
%[3]s
EOS
echo "%[4]s" > "${ResultDestDir}/PWD"

if [ -n "$Hook" ] ; then
  $Hook "${ResultDestDir}"
fi
if [ -n "$Play" ] ; then
  result_wav="${ResultDestDir}/${BASENAME}.wav"
  world_wav="${ResultDestDir}/${BASENAME}_world.wav"
  if [ -f "${result_wav}" ] ; then
    $Play "${result_wav}"
  else
    $Play "${world_wav}"
  fi
fi`,
				g.dir.OutputDir(),
				g.dir.MusicXMLDir(),
				func() []byte {
					b, _ := yaml.Marshal(g.c)
					return b
				}(),
				g.dir.PWD(),
			)))

	return execx.NewExecutableTasks(tasks, g.env())
}
