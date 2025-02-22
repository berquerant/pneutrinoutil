package task

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/pkg/ctl"
	"gopkg.in/yaml.v3"
)

func NewGenerator(dir *Dir, c *ctl.Config, play string) *Generator {
	g := &Generator{
		dir:  dir,
		c:    c,
		play: play,
	}
	return g
}

type Generator struct {
	dir  *Dir
	c    *ctl.Config
	play string
}

func (g Generator) dyldLibraryPath() string {
	x, _ := filepath.Abs(g.dir.NeutrinoDir())
	return fmt.Sprintf(
		"%s/bin:%s",
		x,
		os.Getenv("DYLD_LIBRARY_PATH"),
	)
}

func (g Generator) env() execx.Env {
	e := g.c.Env()
	e.Set("ResultDestDir", g.dir.ResultDestDir())
	e.Set("Score", g.c.Score)
	e.Set("Play", g.play)
	e.Set("DYLD_LIBRARY_PATH", g.dyldLibraryPath())
	e.Set("HOME", os.Getenv("HOME"))
	e.Merge(g.dir.Env())
	return e
}

func (g Generator) ExecutableTasks() *execx.ExecutableTasks {
	tasks := execx.NewTasks().
		Add(execx.NewTask(
			"init",
			fmt.Sprintf(`chmod 755 %[1]s/*
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
				`world_wav="%[1]s/${BASENAME}_world.wav"
cp %[1]s/${BASENAME}.* "%[2]s/${BASENAME}.musicxml" "${world_wav}" "${ResultDestDir}/"
if [ -f "${world_wav}" ] ; then
  cp "${world_wav}" "${ResultDestDir}/"
fi
cat <<EOS > "${ResultDestDir}/config.yml"
%[3]s
EOS
echo "%[4]s" > "${ResultDestDir}/PWD"

ls -lah "${ResultDestDir}/"
if [ -n "$Play" ] ; then
  $Play "${ResultDestDir}/${BASENAME}.wav" || $Play "${world_wav}"
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
