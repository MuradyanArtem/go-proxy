package web

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"proxy/internal/domain/models"
	"proxy/internal/domain/repository"
	"proxy/internal/interfaces/web/server"
	"strconv"

	"github.com/sirupsen/logrus"
)

type ProxyInformation struct {
	InterceptedHttpsRequest *http.Request
	ForwardedHttpsRequest   *http.Request
	Scheme                  string
	Config                  *tls.Config
}

type Sniffer struct {
	request repository.Request
	config  server.Config
}

func NewSniffer(app *repository.Proxy, c *server.Config) *Sniffer {
	return &Sniffer{
		request: app.Request,
		config:  *c,
	}
}

func (s *Sniffer) Recording(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		s.tunnel(w, r)
		return
	}

	s.proxy(w, r)
}

func (s *Sniffer) proxy(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "proxy",
		}).Error(err)
	}

	req := models.Request{
		Request: string(dump),
		URL:     r.RequestURI,
		Headers: r.Header,
	}

	if err = s.request.Insert(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "proxy",
		}).Error(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	defer func() {
		if err = resp.Body.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "web",
				"func": "proxy",
			}).Error(err)
		}
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	copyHeaders(resp.Header, w.Header())

	w.WriteHeader(resp.StatusCode)
	if _, err = io.Copy(w, resp.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "proxy",
		}).Error(err)
	}
}

func fillInformation(r *http.Request) (*ProxyInformation, error) {
	requestedUrl, err := url.Parse(r.RequestURI)
	if err != nil {
		return nil, err
	}

	info := ProxyInformation{}
	if requestedUrl.Scheme == "" {
		info.Scheme = r.URL.Host
	} else {
		info.Scheme = requestedUrl.Scheme
	}

	info.InterceptedHttpsRequest = r
	return &info, nil
}

func hijackConnect(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return nil, errors.New("hijacker !ok")
	}

	hijackedConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return nil, err
	}

	_, err = hijackedConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		hijackedConn.Close()
		return nil, err
	}

	return hijackedConn, nil
}

func generateCertificate(proxyInfo *ProxyInformation) (*tls.Certificate, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cmdGenDir := rootDir + "/ssl/"
	certFilename := cmdGenDir + proxyInfo.Scheme + ".crt"

	if _, errStat := os.Stat(certFilename); os.IsNotExist(errStat) {
		genCommand := exec.Command("sh", cmdGenDir+"/gen_cert.sh", proxyInfo.Scheme, strconv.Itoa(rand.Intn(1000)))
		if _, err := genCommand.CombinedOutput(); err != nil {
			return nil, err
		}
	}

	cert, err := tls.LoadX509KeyPair(certFilename, cmdGenDir+"/cert.key")
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func initializeTCPClient(hijackedConn net.Conn, proxyInfo *ProxyInformation) (*tls.Conn, error) {
	cert, err := generateCertificate(proxyInfo)
	if err != nil {
		return nil, err
	}

	proxyInfo.Config = &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ServerName:   proxyInfo.Scheme,
	}

	TCPClientConn := tls.Server(hijackedConn, proxyInfo.Config)

	if err := TCPClientConn.Handshake(); err != nil {
		TCPClientConn.Close()
		hijackedConn.Close()
		return nil, err
	}

	clientReader := bufio.NewReader(TCPClientConn)
	proxyInfo.ForwardedHttpsRequest, err = http.ReadRequest(clientReader)
	if err != nil {
		return nil, err
	}

	return TCPClientConn, nil
}

func doHttpsRequest(TCPClientConn *tls.Conn, TCPServerConn *tls.Conn, proxyInfo *ProxyInformation) error {
	rawReq, err := httputil.DumpRequest(proxyInfo.ForwardedHttpsRequest, true)
	if err != nil {
		return err
	}

	_, err = TCPServerConn.Write(rawReq)
	if err != nil {
		return err
	}

	serverReader := bufio.NewReader(TCPServerConn)
	TCPServerResponse, err := http.ReadResponse(serverReader, proxyInfo.ForwardedHttpsRequest)
	if err != nil {
		return err
	}

	decodedResponse, err := decodeResponse(TCPServerResponse)
	if err != nil {
		return err
	}

	if _, err = TCPClientConn.Write(decodedResponse); err != nil {
		return err
	}

	return nil
}

func (s *Sniffer) tunnel(w http.ResponseWriter, r *http.Request) {
	hijackedConn, err := hijackConnect(w)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "tunnel",
		}).Error(err)
		return
	}
	defer hijackedConn.Close()

	proxy, err := fillInformation(r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "tunnel",
		}).Error(err)
		return
	}

	TCPClientConn, err := initializeTCPClient(hijackedConn, proxy)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "tunnel",
		}).Error(err)
		return
	}
	defer TCPClientConn.Close()

	TCPServerConn, err := tls.Dial("tcp", proxy.InterceptedHttpsRequest.Host, proxy.Config)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "tunnel",
		}).Error(err)
		return
	}

	err = doHttpsRequest(TCPClientConn, TCPServerConn, proxy)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "tunnel",
		}).Error(err)
		return
	}

	dump, err := httputil.DumpRequest(proxy.ForwardedHttpsRequest, true)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "tunnel",
		}).Error(err)
		return
	}

	err = s.request.Insert(&models.Request{
		Headers: proxy.ForwardedHttpsRequest.Header,
		Request: string(dump),
		URL:     fmt.Sprintf("https://%s%s", proxy.ForwardedHttpsRequest.Host, proxy.ForwardedHttpsRequest.URL.Path),
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "web",
			"func": "tunnel",
		}).Error(err)
		return
	}
}
