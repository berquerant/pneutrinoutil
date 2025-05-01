package repo

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/berquerant/pneutrinoutil/pkg/alog"
	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
)

const (
	ObjectTable = "objects"
)

type CreateObjectRequest struct {
	Type      domain.ObjectType
	Bucket    string
	Path      string
	SizeBytes uint64
}
type ObjectCreator interface {
	CteateObject(ctx context.Context, req *CreateObjectRequest) (*domain.Object, error)
}

type ObjectGetter interface {
	GetObject(ctx context.Context, id int) (*domain.Object, error)
	GetObjectByPath(ctx context.Context, bucket, path string) (*domain.Object, error)
}

func NewObject(query infra.Queryer[domain.Object], exec infra.Execer) *Object {
	return &Object{
		query: query,
		exec:  exec,
	}
}

type Object struct {
	query infra.Queryer[domain.Object]
	exec  infra.Execer
}

var (
	_ ObjectCreator = &Object{}
	_ ObjectGetter  = &Object{}
)

func (*Object) scan(f func(...any) error) (*domain.Object, error) {
	var (
		id        int
		typeId    int
		bucket    string
		path      string
		sizeBytes uint64
		createdAt time.Time
		updatedAt time.Time
		// bucketPathSha256 []byte
	)
	// if err := f(&id, &typeId, &bucket, &path, &sizeBytes, &createdAt, &updatedAt, &bucketPathSha256); err != nil {
	// 	return nil, err
	// }
	if err := f(&id, &typeId, &bucket, &path, &sizeBytes, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	return &domain.Object{
		ID:        id,
		Type:      domain.ObjectType(typeId),
		Bucket:    bucket,
		Path:      path,
		SizeBytes: sizeBytes,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		// BucketPathSha256: bucketPathSha256,
	}, nil
}

func (*Object) bucketPathSha256(bucket, path string) []byte {
	x := sha256.Sum256([]byte(bucket + "___" + path))
	return x[:]
}

func (s *Object) GetObjectByPath(ctx context.Context, bucket, path string) (*domain.Object, error) {
	// hash := s.bucketPathSha256(bucket, path)
	alog.L().Debug("GetObjectByPath", "bucket", bucket, "path", path)
	r, err := s.query.Query(ctx, &infra.QueryRequest[domain.Object]{
		// Query: "select id, type_id, bucket, path, size_bytes, created_at, updated_at, bucket_path_sha256 from objects where bucket_path_sha256 = ? and bucket = ? and path = ?;",
		Query: "select id, type_id, bucket, path, size_bytes, created_at, updated_at from objects where bucket = ? and path = ?;",
		Args: []any{
			// hash,
			bucket,
			path,
		},
		Scan: s.scan,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: get object by path: bucket=%s, path=%s", err, bucket, path)
	}
	if err := r.AssertRows(1); err != nil {
		return nil, fmt.Errorf("%w: get object by path: bucket=%s, path=%s", err, bucket, path)
	}
	return r.Items[0], nil
}

func (s *Object) GetObject(ctx context.Context, id int) (*domain.Object, error) {
	r, err := s.query.Query(ctx, &infra.QueryRequest[domain.Object]{
		// Query: "select id, type_id, bucket, path, size_bytes, created_at, updated_at, bucket_path_sha256 from objects where id = ?;",
		Query: "select id, type_id, bucket, path, size_bytes, created_at, updated_at from objects where id = ?;",
		Args: []any{
			id,
		},
		Scan: s.scan,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: get object: id=%d", err, id)
	}
	if err := r.AssertRows(1); err != nil {
		return nil, fmt.Errorf("%w: get object: id=%d", err, id)
	}
	return r.Items[0], nil
}

func (s *Object) CteateObject(ctx context.Context, req *CreateObjectRequest) (*domain.Object, error) {
	r, err := s.exec.Exec(ctx, &infra.ExecRequest{
		Query: "insert into objects (type_id, bucket, path, size_bytes) values (?, ?, ?, ?);",
		Args: []any{
			int(req.Type),
			req.Bucket,
			req.Path,
			req.SizeBytes,
		},
		AssertResponse: infra.AssertJoin(
			infra.AssertRowsAffected(1),
			infra.AssertLastInserted(),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: create object: bucket=%s, path=%s", err, req.Bucket, req.Path)
	}

	id, _ := r.LastInsertId()
	return s.GetObject(ctx, int(id))
}
