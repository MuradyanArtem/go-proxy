package app

import (
	"proxy/internal/domain/models"
	"proxy/internal/domain/repository"
)

type Request struct {
	db repository.Request
}

func NewRequest(db repository.Request) *Request {
	return &Request{db: db}
}

func (r *Request) Insert(req models.Request) error {
	return r.db.Insert(req)
}

func (r *Request) GetRequestList() ([]models.Request, error) {
	return r.db.GetRequestList()
}

func (r *Request) GetRequestById(id int64) (models.Request, error) {
	return r.db.GetRequestById(id)
}

func (r *Request) GetRequestHeaders(id int64) (models.Request, error) {
	return r.db.GetRequestHeaders(id)
}
