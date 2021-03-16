package web

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"proxy/internal/domain/models"
	"proxy/internal/domain/repository"
	"proxy/internal/interfaces/web/server"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Admin struct {
	request repository.Request
	config  server.Config
}

func NewAdmin(app *repository.Proxy, c *server.Config) *Admin {
	return &Admin{
		request: app.Request,
		config:  *c,
	}
}

func (s *Admin) Init() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/repeat/{id}", s.Repeat).
		Methods(http.MethodGet)

	r.HandleFunc("/requests", s.AllRequests).
		Methods(http.MethodGet)

	r.HandleFunc("/request/{id}", s.RequestById).
		Methods(http.MethodGet)

	r.HandleFunc("/scan/{id}", s.ScanRequest).
		Methods(http.MethodGet)

	return r
}

func createNewRequest(storedRequest *models.Request) (*http.Request, error) {
	requestReader := bufio.NewReader(strings.NewReader(storedRequest.Request))
	buffer, err := http.ReadRequest(requestReader)
	if err != nil {
		return nil, err
	}

	newRequest, err := http.NewRequest(buffer.Method, storedRequest.URL, buffer.Body)
	if err != nil {
		return nil, err
	}

	copyHeaders(buffer.Header, newRequest.Header)
	newRequest.Header.Del("Proxy-Connection")

	return newRequest, nil
}

func (s *Admin) Repeat(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		// log
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 3,
	}
	defer client.CloseIdleConnections()

	storedRequest, err := s.request.GetRequestById(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	newRequest, err := createNewRequest(storedRequest)
	if err != nil {
		// log
		return
	}

	resp, err := client.Do(newRequest)
	if err != nil {
		// log
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		// log
	}
}

func (s *Admin) AllRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	requests, err := s.request.GetRequestList()
	if err != nil {
		// log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(*requests)
	if err != nil {
		// log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		// log
	}
}

func (s *Admin) RequestById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		// log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request, err := s.request.GetRequestById(int64(id))
	if err != nil {
		// log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(*request)
	if err != nil {
		// log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		// log
	}
}

func (s *Admin) ScanRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
