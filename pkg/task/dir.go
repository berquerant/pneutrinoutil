package task

import (
	"path/filepath"

	"github.com/berquerant/execx"
)

func NewDir(workDir, neutrinoDir, pwd string) *Dir {
	return &Dir{
		workDir:     workDir,
		neutrinoDir: neutrinoDir,
		pwd:         pwd,
	}
}

type Dir struct {
	workDir     string
	neutrinoDir string
	pwd         string
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

func (Dir) ResultDir() string { return join("${WORKDIR}", "result") }

func (Dir) ModelDir() string  { return "./model" }
func (Dir) BinDir() string    { return "./bin" }
func (Dir) OutputDir() string { return "./output" }
func (Dir) ScoreDir() string  { return "./score" }

func (d Dir) MusicXMLDir() string { return join(d.ScoreDir(), "musicxml") }
func (d Dir) LabelDir() string    { return join(d.ScoreDir(), "label") }

func (d Dir) FullDir() string   { return join(d.LabelDir(), "full") }
func (d Dir) MonoDir() string   { return join(d.LabelDir(), "mono") }
func (d Dir) TimingDir() string { return join(d.LabelDir(), "timing") }

func join(elem ...string) string { return filepath.Join(elem...) }
