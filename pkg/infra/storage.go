package infra

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/pathx"
	"github.com/berquerant/pneutrinoutil/pkg/ptr"
)

type Object interface {
	ObjectCreator
	ObjectGetter
}

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
	_ Object = &S3{}
)

func NewS3(client *s3.Client) *S3 {
	return &S3{
		client: client,
	}
}

type S3 struct {
	client *s3.Client
}

func (s *S3) CreateObject(ctx context.Context, req *CreateObjectRequest) (*CreateObjectResponse, error) {
	if _, err := req.Object.Blob.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("%w: s3 create object bucket=%s, path=%s", err, req.Object.Bucket, req.Object.Path)
	}
	sizeBytes, err := req.Object.Blob.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("%w: s3 create object bucket=%s, path=%s", err, req.Object.Bucket, req.Object.Path)
	}
	if _, err := req.Object.Blob.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("%w: s3 create object bucket=%s, path=%s", err, req.Object.Bucket, req.Object.Path)
	}

	if _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: ptr.To(req.Object.Bucket),
		Key:    ptr.To(req.Object.Path),
		Body:   req.Object.Blob, // Body should be seekable (implement io.Seeker)
	}); err != nil {
		return nil, fmt.Errorf("%w: s3 create object bucket=%s, path=%s", err, req.Object.Bucket, req.Object.Path)
	}
	return &CreateObjectResponse{
		SizeBytes: uint64(sizeBytes),
	}, nil
}

func (s *S3) GetObject(ctx context.Context, req *GetObjectRequest) (*domain.StorageObject, error) {
	r, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: ptr.To(req.Bucket),
		Key:    ptr.To(req.Path),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: s3 get object bucket=%s, path=%s", err, req.Bucket, req.Path)
	}
	defer r.Body.Close()

	var buf bytes.Buffer
	sizeBytes, err := io.Copy(&buf, r.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: s3 get object bucket=%s, path=%s", err, req.Bucket, req.Path)
	}
	return &domain.StorageObject{
		Bucket:    req.Bucket,
		Path:      req.Path,
		Blob:      bytes.NewReader(buf.Bytes()),
		SizeBytes: uint64(sizeBytes),
	}, nil
}

var (
	_ Object = &FileSystem{}
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
		return nil, fmt.Errorf("%w: fileystsem get object bucket=%s, path=%s", err, req.Bucket, req.Path)
	}
	defer fp.Close()
	sizeBytes, err := io.Copy(&buf, fp)
	if err != nil {
		return nil, fmt.Errorf("%w: filesystem get object bucket=%s, path=%s", err, req.Bucket, req.Path)
	}
	return &domain.StorageObject{
		Bucket:    req.Bucket,
		Path:      req.Path,
		Blob:      bytes.NewReader(buf.Bytes()),
		SizeBytes: uint64(sizeBytes),
	}, nil
}

func (f *FileSystem) CreateObject(_ context.Context, req *CreateObjectRequest) (*CreateObjectResponse, error) {
	path := filepath.Join(f.rootDir, req.Object.Bucket, req.Object.Path)
	if err := pathx.EnsureDir(filepath.Dir(path)); err != nil {
		return nil, fmt.Errorf("%w: filesystem create object bucket=%s, path=%s", err, req.Object.Bucket, req.Object.Path)
	}

	fp, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("%w: filesystem create object bucket=%s, path=%s", err, req.Object.Bucket, req.Object.Path)
	}
	defer fp.Close()

	sizeBytes, err := io.Copy(fp, req.Object.Blob)
	if err != nil {
		return nil, fmt.Errorf("%w: filesystem create object bucket=%s, path=%s", err, req.Object.Bucket, req.Object.Path)
	}
	return &CreateObjectResponse{
		SizeBytes: uint64(sizeBytes),
	}, nil
}
