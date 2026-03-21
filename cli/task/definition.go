package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/cli/ctl"
	"github.com/goccy/go-yaml"
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

func (g Generator) executableTasksV3() *execx.ExecutableTasks {
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
				`%[1]s/neutrino \
  "%[2]s/${BASENAME}.lab" \
  "%[3]s/${BASENAME}.lab" \
  "%[4]s/${BASENAME}.f0" \
  "%[4]s/${BASENAME}.melspec" \
  "%[4]s/${BASENAME}.wav" \
  "%[5]s/${ModelDir}/" %[6]s \
  -n ${NumThreads} \
  -f ${Transpose} \
  -i "%[4]s/${BASENAME}.trace" \
  -t`,
				g.dir.BinDir(),
				g.dir.FullDir(),
				g.dir.TimingDir(),
				g.dir.OutputDir(),
				g.dir.ModelDir(),
				func() string {
					if g.c.SupportModelDir != "" {
						return `-S "` + g.dir.ModelDir() + `/${SupportModelDir}/"`
					}
					return ""
				}(),
			))).
		Add(execx.NewTask(
			"cleanup",
			fmt.Sprintf(
				`cp %[1]s/${BASENAME}.* "%[2]s/${BASENAME}.musicxml" "${ResultDestDir}/"
cat <<EOS > "${ResultDestDir}/config.yml"
%[3]s
EOS
echo "%[4]s" > "${ResultDestDir}/PWD"

if [ -n "$Hook" ] ; then
  $Hook "${ResultDestDir}"
fi
if [ -n "$Play" ] ; then
  result_wav="${ResultDestDir}/${BASENAME}.wav"
  $Play "${result_wav}"
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

func (g Generator) ExecutableTasks() (*execx.ExecutableTasks, error) {
	if strings.Contains(g.c.NeutrinoVersion, "v3.") {
		return g.executableTasksV3(), nil
	}
	return nil, fmt.Errorf("failed to generate executable tasks")
}
