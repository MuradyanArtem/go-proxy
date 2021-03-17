package app

import (
	"net/url"
	"proxy/internal/domain/repository"
)

const XSS = `'"><img src onerror=alert()>`

type Scanner struct {
	db repository.Request
}

func newScanner(db repository.Request) repository.Scanner {
	return &Scanner{
		db: db,
	}
}

func (s *Scanner) MakeInject(req *string) ([]string, error) {
	u, err := url.Parse(*req)
	if err != nil {
		return nil, err
	}

	var res []string
	for param, values := range u.Query() {
		for idx, el := range values {
			bufValues := append(values[:0:0], values...)
			bufValues[idx] = el + XSS

			query := u.Query()
			query[param] = bufValues

			buf := *u
			buf.RawQuery = query.Encode()
			res = append(res, buf.String())
		}
	}
	return res, nil
}
