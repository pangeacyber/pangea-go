package store

import (
	"context"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

type Client interface {
	FolderCreate(ctx context.Context, input *FolderCreateRequest) (*pangea.PangeaResponse[FolderCreateResult], error)
	Delete(ctx context.Context, input *DeleteRequest) (*pangea.PangeaResponse[DeleteResult], error)
	Get(ctx context.Context, input *GetRequest) (*pangea.PangeaResponse[GetResult], error)
	Put(ctx context.Context, input *PutRequest, file *os.File) (*pangea.PangeaResponse[PutResult], error)
	Update(ctx context.Context, input *UpdateRequest) (*pangea.PangeaResponse[UpdateResult], error)
	List(ctx context.Context, input *ListRequest) (*pangea.PangeaResponse[ListResult], error)
	GetArchive(ctx context.Context, input *GetArchiveRequest) (*pangea.PangeaResponse[GetArchiveResult], error)
	ShareLinkCreate(ctx context.Context, input *ShareLinkCreateRequest) (*pangea.PangeaResponse[ShareLinkCreateResult], error)
	ShareLinkGet(ctx context.Context, input *ShareLinkGetRequest) (*pangea.PangeaResponse[ShareLinkGetResult], error)
	ShareLinkList(ctx context.Context, input *ShareLinkListRequest) (*pangea.PangeaResponse[ShareLinkListResult], error)
	ShareLinkDelete(ctx context.Context, input *ShareLinkDeleteRequest) (*pangea.PangeaResponse[ShareLinkDeleteResult], error)
	ShareLinkSend(ctx context.Context, input *ShareLinkSendRequest) (*pangea.PangeaResponse[ShareLinkSendResult], error)
	RequestUploadURL(ctx context.Context, input *PutRequest) (*pangea.PangeaResponse[PutResult], error)

	// Base service methods
	pangea.BaseServicer
}

type store struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &store{
		BaseService: pangea.NewBaseService("store", cfg),
	}
	return cli
}

type FileUploader struct {
	client *pangea.Client
}

func NewFileUploader() FileUploader {
	cfg := &pangea.Config{
		QueuedRetryEnabled: false,
	}

	return FileUploader{
		client: pangea.NewClient("StoreUploader", cfg),
	}
}
