package app

import (
	"proxy/internal/domain/repository"
)

func NewApp(db *repository.Proxy) *repository.Proxy {
	return &repository.Proxy{
		Request: newRequest(db.Request),
	}
}
