package usecase

import (
	"Requester/internal/entity"
	"context"
)

type (
	Record interface {
		Get(ctx context.Context, url string) (entity.Record, bool, error)
		Add(ctx context.Context, r entity.Record) error
		Requester(ctx context.Context, url []string) ([]entity.RecordDto, error)
	}

	RecordRepo interface {
		Set(ctx context.Context, r entity.Record) error
		Get(ctx context.Context, url string) (entity.Record, bool, error)
	}

	RecordWebAPI interface {
		Processing(ctx context.Context, url string) (entity.Record, error)
	}
)
