package pathx_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/stretchr/testify/assert"
)

func TestEnsure(t *testing.T) {
	root := t.TempDir()

	dir := filepath.Join(root, "dir")
	assert.Equal(t, pathx.EnotExist, pathx.Exist(dir))

	assert.Nil(t, pathx.EnsureDir(dir))
	assert.Equal(t, pathx.Edir, pathx.Exist(dir))

	assert.Nil(t, pathx.EnsureDir(dir))
	assert.Equal(t, pathx.Edir, pathx.Exist(dir))

	file := filepath.Join(root, "file")
	const fileMode = 0644
	assert.Equal(t, pathx.EnotExist, pathx.Exist(file))

	assert.Nil(t, pathx.EnsureFile(file, pathx.WithMode(fileMode)))
	assert.Equal(t, pathx.Efile, pathx.Exist(file))
	assertFileMode(t, file, fileMode)

	assert.Nil(t, pathx.EnsureFile(file, pathx.WithMode(fileMode)))
	assert.Equal(t, pathx.Efile, pathx.Exist(file))
	assertFileMode(t, file, fileMode)

	const newFileMode = 0666
	assert.Nil(t, pathx.EnsureFile(file, pathx.WithMode(newFileMode)))
	assert.Equal(t, pathx.Efile, pathx.Exist(file))
	assertFileMode(t, file, newFileMode)

	{
		f, err := os.OpenFile(file, os.O_WRONLY, newFileMode)
		assert.Nil(t, err)
		defer f.Close()
		_, err = fmt.Fprintln(f, "test")
		assert.Nil(t, err)
	}
	assertFileSize := func(t *testing.T, size int64) {
		info, err := os.Stat(file)
		assert.Nil(t, err)
		assert.Equal(t, size, info.Size())
	}
	assertFileSize(t, 5)

	assert.Nil(t, pathx.EnsureFile(file, pathx.WithMode(newFileMode)))
	assertFileSize(t, 5)

	assert.Nil(t, pathx.EnsureFile(file, pathx.WithMode(newFileMode), pathx.WithTruncate(true)))
	assertFileSize(t, 0)
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
