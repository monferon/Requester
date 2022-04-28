package usecase

import (
	"Requester/internal/entity"
	"context"
	"fmt"
	"time"
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

func (ruc *RecordUseCase) Requester(ctx context.Context, records []string) ([]entity.RecordDto, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	guard := make(chan struct{}, 3)
	result := make(chan entity.RecordDto)
	errs := make(chan error, 1)
	resultArr := make([]entity.RecordDto, 0)
	//wg := sync.WaitGroup{}
	for index, val := range records {
		guard <- struct{}{}
		if index > 100 {
			return resultArr, nil
		}
		//wg.Add(1)
		go func(val string) {
			//defer wg.Done()
			record, ok, err := ruc.repo.Get(ctx, val)
			if err != nil {
				errs <- err
			}
			if !ok || time.Now().Unix()-record.Ttl > 10 {
				rec, err := ruc.webAPI.Processing(ctx, val)
				if err != nil {
					errs <- err
				}
				err = ruc.repo.Set(ctx, rec)
				if err != nil {
					errs <- err
				}
				result <- entity.RecordDto{URL: rec.URL, Size: rec.Size}
			} else {
				result <- entity.RecordDto{URL: record.URL, Size: record.Size}
			}
			<-guard
		}(val)
		//wg.Wait()
		select {
		case t := <-result:
			resultArr = append(resultArr, t)
		case e := <-errs:
			//ctx.Done()
			return nil, e

		case <-ctx.Done():
			fmt.Println("ABORTED11111")

		}

	}
	return resultArr, nil
}
