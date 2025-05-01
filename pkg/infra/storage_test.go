package infra_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
	"github.com/stretchr/testify/assert"
)

func TestFileSystem(t *testing.T) {
	fs := infra.NewFileSystem(t.TempDir())

	const (
		bucket  = "test"
		path    = "dir1/obj.txt"
		content = "some_context"
	)

	ctx := context.TODO()
	t.Run("missing", func(t *testing.T) {
		_, err := fs.GetObject(ctx, &infra.GetObjectRequest{
			Bucket: bucket,
			Path:   path,
		})
		assert.NotNil(t, err)
	})
	t.Run("create", func(t *testing.T) {
		buf := bytes.NewBufferString(content)
		_, err := fs.CreateObject(ctx, &infra.CreateObjectRequest{
			Object: &domain.StorageObject{
				Bucket: bucket,
				Path:   path,
				Blob:   buf,
			},
		})
		assert.Nil(t, err)
	})
	t.Run("get", func(t *testing.T) {
		got, err := fs.GetObject(ctx, &infra.GetObjectRequest{
			Bucket: bucket,
			Path:   path,
		})
		if !assert.Nil(t, err) {
			return
		}
		s, err := io.ReadAll(got.Blob)
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, content, string(s))
	})
}
