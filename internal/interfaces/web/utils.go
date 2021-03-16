package web

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

func copyHeaders(from, to http.Header) {
	for header, values := range from {
		for _, value := range values {
			to.Add(header, value)
		}
	}
}

func decodeResponse(response *http.Response) ([]byte, error) {
	var body io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		body, _ = gzip.NewReader(response.Body)
	default:
		body = response.Body
	}

	bodyByte, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	lineBreak := []byte("\n")
	bodyByte = append(bodyByte, lineBreak...)

	var headers string
	for header, values := range response.Header {
		for _, value := range values {
			headers += header + ": " + value + "\n"
		}
	}

	status := response.Status + "\n"
	proto := response.Proto + "\n"
	headers = status + proto + headers

	headersByteArray := []byte(headers)

	defer body.Close()

	return append(headersByteArray, bodyByte...), nil
}
