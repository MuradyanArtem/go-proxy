package models

import "net/http"

type Request struct {
	Id      int64
	Headers http.Header
	Host    string
	Request string
}
