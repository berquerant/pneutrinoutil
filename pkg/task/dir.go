package task

import (
	"fmt"
	"math"
	"math/rand/v2"
	"path/filepath"
	"time"

	"github.com/berquerant/execx"
)

func NewDir(workDir, neutrinoDir, pwd string, now time.Time) *Dir {
	return &Dir{
		workDir:     workDir,
		neutrinoDir: neutrinoDir,
		pwd:         pwd,
		now:         now,
	}
}

type Dir struct {
	workDir     string
	neutrinoDir string
	pwd         string
	now         time.Time
}

func (d Dir) Env() execx.Env {
	e := execx.NewEnv()
	e.Set("WORKDIR", d.workDir)
	// e.Set("NEUTRINODIR", d.neutrinoDir)
	return e
}

func (d Dir) PWD() string         { return d.pwd }
func (d Dir) WorkDir() string     { return d.workDir }
func (d Dir) NeutrinoDir() string { return d.neutrinoDir }

func (d Dir) ResultDir() string { return d.join("${WORKDIR}", "result") }

func (Dir) ModelDir() string  { return "./model" }
func (Dir) BinDir() string    { return "./bin" }
func (Dir) OutputDir() string { return "./output" }
func (Dir) ScoreDir() string  { return "./score" }

func (d Dir) MusicXMLDir() string { return d.join(d.ScoreDir(), "musicxml") }
func (d Dir) LabelDir() string    { return d.join(d.ScoreDir(), "label") }

func (d Dir) FullDir() string   { return d.join(d.LabelDir(), "full") }
func (d Dir) MonoDir() string   { return d.join(d.LabelDir(), "mono") }
func (d Dir) TimingDir() string { return d.join(d.LabelDir(), "timing") }

func (Dir) join(elem ...string) string { return filepath.Join(elem...) }
func (Dir) salt() uint16               { return uint16(rand.IntN(math.MaxUint16 + 1)) }

func (d Dir) timeString() string {
	return fmt.Sprintf(
		// %Y%m%d-%H%M%S
		"%04d%02d%02d-%02d%02d%02d",
		d.now.Year(),
		d.now.Month(),
		d.now.Day(),
		d.now.Hour(),
		d.now.Minute(),
		d.now.Second(),
	)
}

func (d Dir) ResultDestDir() string {
	return d.join(
		d.ResultDir(),
		fmt.Sprintf("${BASENAME}_%s_%d_%d",
			d.timeString(),
			d.now.Unix(),
			d.salt(),
		),
	)
}
