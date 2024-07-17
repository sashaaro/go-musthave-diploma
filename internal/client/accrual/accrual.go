package accrual

import (
	"context"
	"github.com/sashaaro/go-musthave-diploma/internal/client/accrual/dto"
	"github.com/sashaaro/go-musthave-diploma/internal/client/accrual/http/api/orders"
	httpClient "github.com/sashaaro/go-musthave-diploma/internal/client/accrual/http/client"
	"github.com/sashaaro/go-musthave-diploma/pkg/logging"
)

type OrdersAPI interface {
	SendOrder(ctx context.Context, orderDTO dto.Order) (*dto.OrderResponse, error)
}

type client struct {
	host       string
	httpClient httpClient.ClientHTTP
	logger     logging.Logger
	ordersAPI  OrdersAPI
}

func New(host string, logger logging.Logger) *client {
	httpClientInstance := httpClient.New(logger)
	ordersAPI := orders.New(httpClientInstance, host, logger)

	return &client{
		host:       host,
		httpClient: httpClientInstance,
		logger:     logger,
		ordersAPI:  ordersAPI,
	}
}

func (c *client) SendOrder(ctx context.Context, orderDTO dto.Order) (*dto.OrderResponse, error) {
	return c.ordersAPI.SendOrder(ctx, orderDTO)
}
