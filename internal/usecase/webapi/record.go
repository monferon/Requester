package webapi

import (
	"Requester/internal/entity"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RecordWebAPI struct {
}

func New() *RecordWebAPI {
	return &RecordWebAPI{}
}

func (r *RecordWebAPI) Processing(ctx context.Context, url string) (entity.Record, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return entity.Record{}, err
	}
	fmt.Println(url)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.Record{}, err
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return entity.Record{}, err
	}
	record := entity.Record{URL: url, Size: len(f), Ttl: time.Now().Unix()}
	return record, nil
}
