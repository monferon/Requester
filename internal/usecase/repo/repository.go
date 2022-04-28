package repo

import (
	"Requester/internal/entity"
	"context"
	"sync"
)

type RecordRepo struct {
	mu      sync.RWMutex
	Records map[string]entity.Record
}

func New() *RecordRepo {
	return &RecordRepo{
		Records: make(map[string]entity.Record),
	}
}

func (r *RecordRepo) Set(ctx context.Context, record entity.Record) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Records[record.URL] = entity.Record{Size: record.Size, Ttl: record.Ttl, URL: record.URL}
	return nil
}

func (r *RecordRepo) Get(ctx context.Context, URL string) (entity.Record, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	val, ok := r.Records[URL]
	return val, ok, nil
}
