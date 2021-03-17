package web

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"proxy/internal/app"
	"proxy/internal/domain/models"
	"proxy/internal/domain/repository"
	"proxy/internal/interfaces/web/server"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Admin struct {
	request repository.Request
	scanner repository.Scanner
	config  server.Config
}

func NewAdmin(app *repository.Proxy, c *server.Config) *Admin {
	return &Admin{
		request: app.Request,
		scanner: app.Scanner,
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
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "Repeat",
		}).Error(err)
	}

	storedRequest, err := s.request.GetRequestById(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	newRequest, err := createNewRequest(storedRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "Repeat",
		}).Error(err)
		return
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Duration(3 * time.Second),
	}
	defer client.CloseIdleConnections()

	resp, err := client.Do(newRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "Repeat",
		}).Error(err)
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "Repeat",
		}).Error(err)
	}
}

func (s *Admin) AllRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	requests, err := s.request.GetRequestList()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "AllRequests",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(requests)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "AllRequests",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "AllRequests",
		}).Error(err)
	}
}

func (s *Admin) RequestById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "RequestById",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request, err := s.request.GetRequestById(int64(id))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "RequestById",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(*request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "RequestById",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "RequestById",
		}).Error(err)
	}
}

func CheckXSS(response *http.Response) (bool, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	isVulBody := strings.Contains(string(body), app.XSS)
	isVulRequest := strings.Contains(string(response.Request.URL.RawQuery), app.XSS)
	if isVulBody || isVulRequest {
		return true, nil
	}
	return false, nil
}

func (s *Admin) ScanRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "ScanRequest",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req, err := s.request.GetRequestById(int64(id))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "ScanRequest",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	injectedURLs, err := s.scanner.MakeInject(&req.URL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "ScanRequest",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Duration(3 * time.Second),
	}
	defer client.CloseIdleConnections()

	scanner := models.Scanner{}
	for _, url := range injectedURLs {
		injectedReq := *req
		injectedReq.URL = url
		newRequest, err := createNewRequest(&injectedReq)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "web",
				"func": "Repeat",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		injectedResp, err := client.Do(newRequest)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "web",
				"func": "Repeat",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer injectedResp.Body.Close()

		isVulnerable, err := CheckXSS(injectedResp)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "web",
				"func": "Repeat",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if isVulnerable {
			scanner.VulnerableURL = append(scanner.VulnerableURL, url)
		}
	}

	response, err := json.Marshal(scanner)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "Repeat",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(response)
	w.WriteHeader(http.StatusOK)
}
