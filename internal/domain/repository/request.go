package repository

import "proxy/internal/domain/models"

type Request interface {
	Insert(*models.Request) error
	GetRequestList() ([]models.Request, error)
	GetRequestById(int64) (*models.Request, error)
	GetRequestHeaders(int64) (*models.Request, error)
}
