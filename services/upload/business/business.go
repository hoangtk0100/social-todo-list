package business

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
)

type UploadRepository interface {
	CreateImage(ctx context.Context, data *core.Image) error
}

type business struct {
	repo     UploadRepository
	provider core.StorageComponent
}

func NewBusiness(repo UploadRepository, provider core.StorageComponent) *business {
	return &business{
		repo:     repo,
		provider: provider,
	}
}
