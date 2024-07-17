package orders

import (
	"github.com/sashaaro/go-musthave-diploma/internal/client/accrual/http/client"
	"github.com/sashaaro/go-musthave-diploma/pkg/logging"
)

type update struct {
	HTTPClient client.ClientHTTP
	BaseURL    string
	logger     logging.Logger
}

func New(HTTPClient client.ClientHTTP, BaseURL string, logger logging.Logger) *update {
	return &update{
		HTTPClient: HTTPClient,
		BaseURL:    BaseURL,
		logger:     logger,
	}
}
