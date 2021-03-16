package storage

import (
	"proxy/internal/domain/repository"
	"time"

	"github.com/jackc/pgx"
)

type DBConfig struct {
	Host                 string
	Port                 uint16
	User                 string
	Database             string
	Password             string
	PreferSimpleProtocol bool
	MaxConnections       int
	AcquireTimeout       time.Duration
}

func NewDB(conn *pgx.ConnPool) (*repository.Proxy, error) {
	return &repository.Proxy{
		Request: newRequest(conn),
	}, nil
}
