package client

import (
	"github.com/sashaaro/go-musthave-diploma/pkg/logging"
	"io"
	netHTTP "net/http"
)

type HttpClient struct {
	logger logging.Logger
}

func New(logger logging.Logger) *HttpClient {
	return &HttpClient{
		logger: logger,
	}
}

func (h HttpClient) NewRequest(method, url string, body io.Reader) (*netHTTP.Request, error) {
	r, err := netHTTP.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return r, err
}

func (h HttpClient) Do(req *netHTTP.Request) (*netHTTP.Response, error) {
	return netHTTP.DefaultClient.Do(req)
}
