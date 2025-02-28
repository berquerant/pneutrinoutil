package pworker_test

import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/pneutrinoutil/server/pworker"
	"github.com/stretchr/testify/assert"
)

type dirEntry struct {
	name  string
	isDir bool
}

func (e dirEntry) Name() string             { return e.name }
func (e dirEntry) IsDir() bool              { return e.isDir }
func (dirEntry) Type() fs.FileMode          { return fs.ModeIrregular }
func (dirEntry) Info() (fs.FileInfo, error) { return nil, errors.New("NotImplemented") }

func newDirEntry(name string, isDir bool) *dirEntry {
	return &dirEntry{
		name:  name,
		isDir: isDir,
	}
}

func TestLoadResultElements(t *testing.T) {
	for _, tc := range []struct {
		title string
		entry os.DirEntry
		want  *pathx.ResultElement
		err   error
	}{
		{
			title: "not dir",
			entry: newDirEntry("notdir", false),
			err:   pworker.ErrNotDir,
		},
		{
			title: "parse error",
			entry: newDirEntry("parse_error", true),
			err:   pathx.ErrParseResultElement,
		},
		{
			title: "loaded",
			entry: newDirEntry("name__20250227010203_1740585723_100", true),
			want: pathx.NewResultElement(
				"name",
				time.Date(2025, time.February, 27, 1, 2, 3, 0, time.UTC),
				int64(1740585723),
				100,
			),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := pworker.LoadResultElement(tc.entry)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
