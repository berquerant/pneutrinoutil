package infra_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	alog.Setup(os.Stdout, slog.LevelDebug)

	dropBucket := func() error {
		return exec.Command("../../bin/ddl.sh", "drop", "storage", bucket).Run()
	}
	createBucket := func() error {
		return exec.Command("../../bin/ddl.sh", "storage", bucket).Run()
	}
	setUp := func(t *testing.T) bool {
		t.Helper()
		return assert.Nil(t, dropBucket(), "drop") && assert.Nil(t, createBucket(), "create")
	}
	tearDown := func(t *testing.T) bool {
		t.Helper()
		return assert.Nil(t, dropBucket(), "drop")
	}

	if !setUp(t) {
		return
	}
	defer tearDown(t)

	for _, tc := range []struct {
		title string
		param *infra.StorageParam
	}{
		{
			title: "S3",
			param: &infra.StorageParam{
				UseS3: true,
				Debug: true,
			},
		},
		{
			title: "FileSystem",
			param: &infra.StorageParam{
				RootDir: t.TempDir(),
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			s, err := infra.NewStorage(context.TODO(), tc.param)
			if !assert.Nil(t, err) {
				return
			}
			testObject(t, context.TODO(), s)
		})
	}
}

const (
	bucket  = "test"
	path    = "dir1/obj.txt"
	content = "some_context"
)

func testObject(t *testing.T, ctx context.Context, s infra.Object) {
	t.Run("missing", func(t *testing.T) {
		_, err := s.GetObject(ctx, &infra.GetObjectRequest{
			Bucket: bucket,
			Path:   path,
		})
		assert.NotNil(t, err)
	})
	t.Run("create", func(t *testing.T) {
		buf := bytes.NewBufferString(content)
		_, err := s.CreateObject(ctx, &infra.CreateObjectRequest{
			Object: &domain.StorageObject{
				Bucket: bucket,
				Path:   path,
				Blob:   buf,
			},
		})
		assert.Nil(t, err)
	})
	t.Run("get", func(t *testing.T) {
		got, err := s.GetObject(ctx, &infra.GetObjectRequest{
			Bucket: bucket,
			Path:   path,
		})
		if !assert.Nil(t, err) {
			return
		}
		buf, err := io.ReadAll(got.Blob)
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, content, string(buf))
	})
}
