package webapi

import "Requester/internal/entity"

type RecordWebAPI struct {
}

func New() *RecordWebAPI {
	return &RecordWebAPI{}
}

func (r *RecordWebAPI) Processing([]entity.Record) ([]entity.Record, error) {
	return nil, nil
	//result, err :=
}
