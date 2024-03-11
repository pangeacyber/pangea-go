package pangea

import (
	"context"
	"fmt"
)

type FileUploader struct {
	client *Client
}

func NewFileUploader() FileUploader {
	cfg := &Config{
		QueuedRetryEnabled: false,
	}

	return FileUploader{
		client: NewClient("FileUploader", cfg),
	}
}

func (fu *FileUploader) UploadFile(ctx context.Context, url string, tm TransferMethod, fd FileData) error {
	if tm == TMmultipart {
		return fmt.Errorf("%s is not supported in UploadFile. Use service client instead", tm)
	}

	fds := FileData{
		File:    fd.File,
		Name:    "file",
		Details: fd.Details,
	}
	return fu.client.UploadFile(ctx, url, tm, fds)
}
