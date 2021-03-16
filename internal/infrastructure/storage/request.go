package storage

import (
	"net/http"
	"proxy/internal/domain/models"
	"proxy/internal/domain/repository"

	"github.com/jackc/pgx"
)

type Request struct {
	connection *pgx.ConnPool
}

func newRequest(conn *pgx.ConnPool) repository.Request {
	return &Request{
		connection: conn,
	}
}

func (db *Request) Insert(req *models.Request) error {
	tx, err := db.connection.Begin()
	if err != nil {
		return err
	}

	defer func(tx *pgx.Tx) {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}(tx)

	var id int64
	err = tx.QueryRow("insert into requests (host, request) values ($1, $2) returning id", req.URL, req.Request).
		Scan(&id)

	for key, vval := range req.Headers {
		for _, val := range vval {
			if key == "Proxy-Connection" {
				continue
			}
			_, err = tx.Exec("insert into headers (req_id, key, val) values ($1, $2, $3)", id, key, val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *Request) GetRequestList() (*[]models.Request, error) {
	tx, err := db.connection.Begin()
	if err != nil {
		return nil, err
	}

	defer func(tx *pgx.Tx) {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}(tx)

	out, err := tx.Query("select id, host, request from requests")
	if err != nil {
		return nil, err
	}

	requests := make([]models.Request, 0, 0)
	for out.Next() {
		var request models.Request

		if err = out.Scan(&request.Id, &request.URL, &request.Request); err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	if err = out.Err(); err != nil {
		return nil, err
	}
	return &requests, nil
}

func (db *Request) GetRequestById(id int64) (*models.Request, error) {
	tx, err := db.connection.Begin()
	if err != nil {
		return nil, err
	}

	defer func(tx *pgx.Tx) {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}(tx)

	out := tx.QueryRow("select request, host from requests where id=$1", id)

	request := models.Request{Id: id}
	if err = out.Scan(&request.Request, &request.URL); err != nil {
		return nil, err
	}

	return &request, nil
}

func (db *Request) GetRequestHeaders(id int64) (*models.Request, error) {
	tx, err := db.connection.Begin()
	if err != nil {
		return nil, err
	}

	defer func(tx *pgx.Tx) {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}(tx)

	out, err := tx.Query("select key, val from headers where req_id=$1", id)
	if err != nil {
		return nil, err
	}

	request := models.Request{Id: id, Headers: http.Header{}}
	for out.Next() {
		var key string
		var val string

		if err = out.Scan(&key, &val); err != nil {
			return nil, err
		}

		request.Headers.Add(key, val)
	}

	if err = out.Err(); err != nil {
		return nil, err
	}
	return &request, nil
}
