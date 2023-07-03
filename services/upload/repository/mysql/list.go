package mysql

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (repo *mysqlRepo) ListImages(
	ctx context.Context,
	ids []int,
	moreKeys ...string,
) ([]core.Image, error) {
	_, span := trace.StartSpan(ctx, "upload.repository.mysql.list")
	defer span.End()

	var result []core.Image

	if err := repo.db.
		Table(core.Image{}.TableName()).
		Where("id in ?", ids).
		Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
