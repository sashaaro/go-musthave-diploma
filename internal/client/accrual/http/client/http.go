package client

import (
	"github.com/sashaaro/go-musthave-diploma-tpl/pkg/logging"
	"io"
	netHTTP "net/http"
	"time"
)

type httpClient struct {
	HashKey *string
	logger  logging.Logger
}

func New(logger logging.Logger) ClientHTTP {
	return &httpClient{
		logger: logger,
	}
}

func (h httpClient) NewRequest(method, url string, body io.Reader) (*netHTTP.Request, error) {
	time.Sleep(time.Second)
	r, err := netHTTP.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return r, err
}

func (h httpClient) Do(req *netHTTP.Request) (*netHTTP.Response, error) {
	return netHTTP.DefaultClient.Do(req)
}
