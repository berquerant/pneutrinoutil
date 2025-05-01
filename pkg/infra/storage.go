package infra

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
)

type CreateObjectRequest struct {
	Object *domain.StorageObject
}

type CreateObjectResponse struct {
	SizeBytes uint64
}

type ObjectCreator interface {
	CreateObject(ctx context.Context, req *CreateObjectRequest) (*CreateObjectResponse, error)
}

type GetObjectRequest struct {
	Bucket string
	Path   string
}

type ObjectGetter interface {
	GetObject(ctx context.Context, req *GetObjectRequest) (*domain.StorageObject, error)
}

var (
	_ ObjectCreator = &FileSystem{}
	_ ObjectGetter  = &FileSystem{}
)

func NewFileSystem(rootDir string) *FileSystem {
	return &FileSystem{
		rootDir: rootDir,
	}
}

type FileSystem struct {
	rootDir string
}

func (f *FileSystem) GetObject(_ context.Context, req *GetObjectRequest) (*domain.StorageObject, error) {
	path := filepath.Join(f.rootDir, req.Bucket, req.Path)
	var buf bytes.Buffer
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	sizeBytes, err := io.Copy(&buf, fp)
	if err != nil {
		return nil, err
	}
	return &domain.StorageObject{
		Bucket:    req.Bucket,
		Path:      req.Path,
		Blob:      &buf,
		SizeBytes: uint64(sizeBytes),
	}, nil
}

func (f *FileSystem) CreateObject(_ context.Context, req *CreateObjectRequest) (*CreateObjectResponse, error) {
	path := filepath.Join(f.rootDir, req.Object.Bucket, req.Object.Path)
	if err := pathx.EnsureDir(filepath.Dir(path)); err != nil {
		return nil, err
	}

	fp, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	sizeBytes, err := io.Copy(fp, req.Object.Blob)
	if err != nil {
		return nil, err
	}
	return &CreateObjectResponse{
		SizeBytes: uint64(sizeBytes),
	}, nil
}
