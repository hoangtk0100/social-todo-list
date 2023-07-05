package mysql

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

func (repo *mysqlRepo) DeleteImages(ctx context.Context, ids []int) error {
	_, span := trace.StartSpan(ctx, "upload.repository.mysql.delete")
	defer span.End()

	if err := repo.db.Table(core.Image{}.TableName()).
		Where("id in ?", ids).
		Delete(nil).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
