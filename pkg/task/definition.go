package task

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/berquerant/execx"
	"github.com/berquerant/pneutrinoutil/pkg/ctl"
	"github.com/berquerant/pneutrinoutil/pkg/script"
	"gopkg.in/yaml.v3"
)

func NewGenerator(dir *Dir, c *ctl.Config, now time.Time, play bool) *Generator {
	g := &Generator{
		dir:  dir,
		c:    c,
		now:  now,
		play: play,
	}
	g.resultDir = g.resultDestDir()
	return g
}

type Generator struct {
	dir       *Dir
	c         *ctl.Config
	now       time.Time
	play      bool
	resultDir string
}

func (g Generator) salt() uint16 {
	return uint16(rand.IntN(math.MaxUint16 + 1))
}

func (g Generator) timeString() string {
	return fmt.Sprintf(
		// %Y%m%d-%H%M%S
		"%04d%02d%02d-%02d%02d%02d",
		g.now.Year(),
		g.now.Month(),
		g.now.Day(),
		g.now.Hour(),
		g.now.Minute(),
		g.now.Second(),
	)
}

func (g Generator) resultDestDir() string {
	return filepath.Join(
		g.dir.ResultDir(),
		fmt.Sprintf(
			"${BASENAME}_%d_%s_%d",
			g.now.Unix(),
			g.timeString(),
			g.salt(),
		),
	)
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
	e.Set("ResultDestDir", g.resultDir)
	e.Set("Score", g.c.Score)
	if g.play {
		e.Set("Play", "1")
	} else {
		e.Set("Play", "")
	}
	e.Set("DYLD_LIBRARY_PATH", g.dyldLibraryPath())
	e.Merge(g.dir.Env())
	return e
}

func (g Generator) scriptOpts() []script.ConfigOption {
	return []script.ConfigOption{
		script.WithDir(g.dir.NeutrinoDir()),
		script.WithEnv(g.env()),
	}
}

func (g Generator) newTask(s *script.Script) *Task {
	return &Task{
		s:    s,
		opts: g.scriptOpts(),
	}
}

func (g Generator) DisplayEnv() *Task {
	var (
		xs  []string
		add = func(x string) { xs = append(xs, x) }
		src = g.env().IntoSlice()
	)

	sort.Strings(src)

	add(fmt.Sprintf("# PWD=%s", g.dir.PWD()))
	for _, x := range src {
		ss := strings.SplitN(x, "=", 2)
		k, v := ss[0], ss[1]
		add(fmt.Sprintf(`%s="%s"`, k, v))
	}
	add(fmt.Sprintf(`export DYLD_LIBRARY_PATH="%s"`, g.dyldLibraryPath()))

	return g.newTask(script.New(
		"display_env",
		strings.Join(xs, "\n"),
	))
}

func (g Generator) Init() *Task {
	return g.newTask(script.New(
		"init",
		fmt.Sprintf(`chmod 755 %[1]s/*
xattr -dr com.apple.quarantine "%[1]s"`,
			g.dir.BinDir(),
		),
	))
}

func (g Generator) Prepare() *Task {
	return g.newTask(script.New(
		"prepare",
		fmt.Sprintf(`cp -f %s "%s/"
mkdir -p ${ResultDestDir}`,
			g.c.Score, g.dir.MusicXMLDir(),
		),
	))
}

func (g Generator) MusicXMLToLabel() *Task {
	return g.newTask(script.New(
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
		),
	))
}

func (g Generator) NEUTRINO() *Task {
	return g.newTask(script.New(
		"NEUTRINO",
		fmt.Sprintf(
			`%s/NEUTRINO \
  "%s/${BASENAME}.lab" \
  "%s/${BASENAME}.lab" \
  "%s/${BASENAME}.f0" \
  "%s/${BASENAME}.melspec" \
  "%s/${ModelDir}/" \
  -w "%s/${BASENAME}.mgc" \
  "%s/${BASENAME}.bap" \
  -n 1 \
  -o ${NumThreads} \
  -k ${StyleShift} \
  -d ${InferenceMode} \
  -t`,
			g.dir.BinDir(),
			g.dir.FullDir(),
			g.dir.TimingDir(),
			g.dir.OutputDir(),
			g.dir.OutputDir(),
			g.dir.ModelDir(),
			g.dir.OutputDir(),
			g.dir.OutputDir(),
		),
	))
}

func (g Generator) NSF() *Task {
	return g.newTask(script.New(
		"NSF",
		fmt.Sprintf(
			`%s/NSF \
  "%s/${BASENAME}.f0" \
  "%s/${BASENAME}.melspec" \
  "%s/${ModelDir}/${NsfModel}.bin" \
  "%s/${BASENAME}.wav" \
  -l "%s/${BASENAME}.lab" \
  -n 1 \
  -p ${NumThreads} \
  -s ${SamplingFreq} \
  -f ${PitchShiftNsf} \
  -t`,
			g.dir.BinDir(),
			g.dir.OutputDir(),
			g.dir.OutputDir(),
			g.dir.ModelDir(),
			g.dir.OutputDir(),
			g.dir.TimingDir(),
		),
	))
}

func (g Generator) WORLD() *Task {
	return g.newTask(script.New(
		"WORLD",
		fmt.Sprintf(
			`%s/WORLD \
  "%s/${BASENAME}.f0" \
  "%s/${BASENAME}.mgc" \
  "%s/${BASENAME}.bap" \
  "%s/${BASENAME}_world.wav" \
  -f ${PitchShiftWorld} \
  -m ${FormantShift} \
  -p ${SmoothPicth} \
  -c ${SmoothFormant} \
  -b ${EnhanceBreathiness} \
  -n ${NumThreads} \
  -t`,
			g.dir.BinDir(),
			g.dir.OutputDir(),
			g.dir.OutputDir(),
			g.dir.OutputDir(),
			g.dir.OutputDir(),
		),
	))
}

func (g Generator) Cleanup() *Task {
	b, _ := yaml.Marshal(g.c)
	return g.newTask(script.New(
		"cleanup",
		fmt.Sprintf(
			`cp "%s/${BASENAME}.wav" \
  "%s/${BASENAME}.musicxml" \
  "${ResultDestDir}/"
cat <<EOS > "${ResultDestDir}/config.yml"
%s
EOS
echo "%s" > "${ResultDestDir}/PWD"

if [ -n "$Play" ] ; then
  open "${ResultDestDir}/${BASENAME}.wav"
fi
`,
			g.dir.OutputDir(),
			g.dir.MusicXMLDir(),
			b,
			g.dir.PWD(),
		),
	))
}
