package repo

import (
	"context"
	"errors"
	"io"

	"github.com/berquerant/pneutrinoutil/pkg/domain"
	"github.com/berquerant/pneutrinoutil/pkg/infra"
)

type ReadObjectResponse interface {
	Object() *domain.Object
	Storage() (*domain.StorageObject, bool)
}

type ObjectReader interface {
	ReadObject(ctx context.Context, id int) (ReadObjectResponse, error)
}

type WriteObjectRequest struct {
	Type   domain.ObjectType
	Bucket string
	Path   string
	Blob   io.Reader
}

type WriteObjectResponse interface {
	Object() *domain.Object
}

type ObjectWriter interface {
	WriteObject(ctx context.Context, req *WriteObjectRequest) (WriteObjectResponse, error)
}

type objectPair struct {
	obj  *domain.Object
	stor *domain.StorageObject
}

func (p *objectPair) Object() *domain.Object                 { return p.obj }
func (p *objectPair) Storage() (*domain.StorageObject, bool) { return p.stor, p.stor != nil }

var (
	_ ReadObjectResponse  = &objectPair{}
	_ WriteObjectResponse = &objectPair{}
	_ ObjectReader        = &ObjectAdmin{}
	_ ObjectWriter        = &ObjectAdmin{}
)

func NewObjectAdmin(
	getter ObjectGetter,
	creator ObjectCreator,
	storageGetter infra.ObjectGetter,
	storageCreator infra.ObjectCreator,
) *ObjectAdmin {
	return &ObjectAdmin{
		creator:        creator,
		getter:         getter,
		storageCreator: storageCreator,
		storageGetter:  storageGetter,
	}
}

type ObjectAdmin struct {
	creator        ObjectCreator
	getter         ObjectGetter
	storageCreator infra.ObjectCreator
	storageGetter  infra.ObjectGetter
}

var (
	ErrUnknownObjectType = errors.New("UnknownObjectType")
)

func (a *ObjectAdmin) ReadObject(ctx context.Context, id int) (ReadObjectResponse, error) {
	obj, err := a.getter.GetObject(ctx, id)
	if err != nil {
		return nil, err
	}
	switch obj.Type {
	case domain.ObjectTypeFile:
		stor, err := a.storageGetter.GetObject(ctx, &infra.GetObjectRequest{
			Bucket: obj.Bucket,
			Path:   obj.Path,
		})
		if err != nil {
			return nil, err
		}
		return &objectPair{
			obj:  obj,
			stor: stor,
		}, nil
	case domain.ObjectTypeDir:
		return &objectPair{
			obj: obj,
		}, nil
	default:
		return nil, ErrUnknownObjectType
	}
}

func (a *ObjectAdmin) WriteObject(ctx context.Context, req *WriteObjectRequest) (WriteObjectResponse, error) {
	switch req.Type {
	case domain.ObjectTypeFile:
		resp, err := a.storageCreator.CreateObject(ctx, &infra.CreateObjectRequest{
			Object: &domain.StorageObject{
				Bucket: req.Bucket,
				Path:   req.Path,
				Blob:   req.Blob,
			},
		})
		if err != nil {
			return nil, err
		}
		obj, err := a.creator.CteateObject(ctx, &CreateObjectRequest{
			Type:      req.Type,
			Bucket:    req.Bucket,
			Path:      req.Path,
			SizeBytes: resp.SizeBytes,
		})
		if err != nil {
			return nil, err
		}
		return &objectPair{
			obj: obj,
		}, nil
	case domain.ObjectTypeDir:
		obj, err := a.createObject(ctx, req)
		if err != nil {
			return nil, err
		}
		return &objectPair{
			obj: obj,
		}, nil
	default:
		return nil, ErrUnknownObjectType
	}
}

func (a *ObjectAdmin) createObject(ctx context.Context, req *WriteObjectRequest) (*domain.Object, error) {
	return a.creator.CteateObject(ctx, &CreateObjectRequest{
		Type:   req.Type,
		Bucket: req.Bucket,
		Path:   req.Path,
	})
}
