package pathx_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/stretchr/testify/assert"
)

func TestEnsure(t *testing.T) {
	t.Run("EnsureDir", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "dir")

		for _, tc := range []struct {
			title    string
			action   func() error
			wantType pathx.ExistType
		}{
			{
				title:    "initial state is not exist",
				action:   func() error { return nil },
				wantType: pathx.EnotExist,
			},
			{
				title:    "create dir",
				action:   func() error { return pathx.EnsureDir(dir) },
				wantType: pathx.Edir,
			},
			{
				title:    "create existing dir is idempotent",
				action:   func() error { return pathx.EnsureDir(dir) },
				wantType: pathx.Edir,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				assert.Nil(t, tc.action())
				assert.Equal(t, tc.wantType, pathx.Exist(dir))
			})
		}
	})

	t.Run("EnsureFile", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "file")

		const initialMode = 0644
		const updatedMode = 0666

		for _, tc := range []struct {
			title        string
			opt          []pathx.ConfigOption
			writeContent string
			wantType     pathx.ExistType
			wantMode     os.FileMode
			wantSize     int64
		}{
			{
				title:    "create new file",
				opt:      []pathx.ConfigOption{pathx.WithMode(initialMode)},
				wantType: pathx.Efile,
				wantMode: initialMode,
				wantSize: 0,
			},
			{
				title:    "ensure existing file without changes",
				opt:      []pathx.ConfigOption{pathx.WithMode(initialMode)},
				wantType: pathx.Efile,
				wantMode: initialMode,
				wantSize: 0,
			},
			{
				title:        "update mode and write content",
				opt:          []pathx.ConfigOption{pathx.WithMode(updatedMode)},
				writeContent: "test\n",
				wantType:     pathx.Efile,
				wantMode:     updatedMode,
				wantSize:     5,
			},
			{
				title:    "ensure file preserves content when truncate is false",
				opt:      []pathx.ConfigOption{pathx.WithMode(updatedMode)},
				wantType: pathx.Efile,
				wantMode: updatedMode,
				wantSize: 5,
			},
			{
				title:    "ensure file truncates when truncate is true",
				opt:      []pathx.ConfigOption{pathx.WithMode(updatedMode), pathx.WithTruncate(true)},
				wantType: pathx.Efile,
				wantMode: updatedMode,
				wantSize: 0,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				assert.Nil(t, pathx.EnsureFile(file, tc.opt...))
				assert.Equal(t, tc.wantType, pathx.Exist(file))
				assertFileMode(t, file, tc.wantMode)

				if tc.writeContent != "" {
					assert.Nil(t, os.WriteFile(file, []byte(tc.writeContent), tc.wantMode))
				}

				info, err := os.Stat(file)
				assert.Nil(t, err)
				assert.Equal(t, tc.wantSize, info.Size())
			})
		}
	})
}

func assertFileMode(t *testing.T, path string, want os.FileMode) {
	info, err := os.Stat(path)
	assert.Nil(t, err)
	assert.Equal(t, want, info.Mode())
}

func TestExist(t *testing.T) {
	root := t.TempDir()

	t.Run("dir", func(t *testing.T) {
		assert.Equal(t, pathx.Edir, pathx.Exist(root))
	})
	t.Run("not exist", func(t *testing.T) {
		assert.Equal(t, pathx.EnotExist, pathx.Exist(filepath.Join(root, "notExist")))
	})
	t.Run("file", func(t *testing.T) {
		pwd, err := os.Getwd()
		assert.Nil(t, err)
		assert.Equal(t, pathx.Efile, pathx.Exist(filepath.Join(pwd, "path_test.go")))
	})
}

func TestBasename(t *testing.T) {
	for _, tc := range []struct {
		title string
		path  string
		want  string
	}{
		{
			title: "no dots",
			path:  "filename",
			want:  "filename",
		},
		{
			title: "one dots",
			path:  "filename.log",
			want:  "filename",
		},
		{
			title: "two dots",
			path:  "filename.log.gz",
			want:  "filename.log",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, pathx.Basename(tc.path))
		})
	}
}
