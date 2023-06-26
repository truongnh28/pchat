package repositories

import (
	"context"
	"database/sql"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

var (
	ErrInvalidParameter = gorm.ErrInvalidData
	ErrNotFound         = gorm.ErrRecordNotFound
)

type Pagination struct {
	Size   int
	Number int
}

func (p Pagination) GetQueryPagingValue() (limit, offset int, ok bool) {
	if p.Size <= 0 || p.Number < 0 || (p.Size == 0 && p.Number == 0) {
		return 0, 0, false
	}

	return p.Size, p.Size * p.Number, true
}

type TimeRange struct {
	StartTime time.Time
	EndTime   time.Time
}

func finalizeTX(tx *sql.Tx, err error) error {
	if err != nil {
		return multierror.Append(err, tx.Rollback())
	}
	return tx.Commit()
}

func finalizeTransaction(tx *gorm.DB, err error) error {
	if err == nil {
		return tx.Commit().Error
	}

	if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
		return multierror.Append(err, rollbackErr)
	}

	return err
}

func fetchPageByTimeRange[T any](
	ctx context.Context,
	paging Pagination,
	timeRange TimeRange,
	db *gorm.DB,
	destTotal *int64,
	destModels *[]T,
) error {
	limit, offset, ok := paging.GetQueryPagingValue()
	if !ok {
		return ErrInvalidParameter
	}

	g, grCtx := errgroup.WithContext(ctx)
	q := db.Model(new(T)).
		Where("created_at BETWEEN ? AND ?", timeRange.StartTime, timeRange.EndTime).
		WithContext(grCtx)
	g.Go(func() error {
		return q.Count(destTotal).Error
	})
	g.Go(func() error {
		return q.Order("created_at DESC").Limit(limit).Offset(offset).Find(destModels).Error
	})

	return g.Wait()
}
