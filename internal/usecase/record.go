package usecase

import (
	"Requester/internal/entity"
	"context"
)

type RecordUseCase struct {
	repo   RecordRepo
	webAPI RecordWebAPI
}

func New(r RecordRepo, w RecordWebAPI) *RecordUseCase {
	return &RecordUseCase{
		repo:   r,
		webAPI: w,
	}
}

func (ruc *RecordUseCase) Get(ctx context.Context, url string) (entity.Record, bool, error) {
	record, ok, err := ruc.repo.Get(ctx, url)

	if err != nil {
		return entity.Record{}, false, err
	}
	return record, ok, nil
}
func (ruc *RecordUseCase) Add(ctx context.Context, record entity.Record) error {
	err := ruc.repo.Set(ctx, record)
	if err != nil {
		return err
	}
	return nil
}
