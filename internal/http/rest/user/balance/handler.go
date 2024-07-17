package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GTech1256/go-musthave-diploma-tpl/internal/domain/entity"
	http2 "github.com/GTech1256/go-musthave-diploma-tpl/internal/http"
	privateRouter "github.com/GTech1256/go-musthave-diploma-tpl/internal/http/middlware/private_router"
	"github.com/GTech1256/go-musthave-diploma-tpl/internal/http/utils/auth"
	"github.com/GTech1256/go-musthave-diploma-tpl/internal/service/user"
	logging2 "github.com/GTech1256/go-musthave-diploma-tpl/pkg/logging"
	"github.com/GTech1256/go-musthave-diploma-tpl/pkg/luhn"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"time"
)

type JWTClient interface {
	BuildJWTString(userId int) (string, error)
	GetTokenExp() time.Duration
	GetUserID(tokenString string) (int, error)
}

type UserExister interface {
	GetIsUserExistById(ctx context.Context, userId int) (bool, error)
}

type Service interface {
	Login(ctx context.Context, userRegister *entity.UserLoginJSON) (*entity.UserDB, error)
	GetById(ctx context.Context, userId int) (*entity.UserDB, error)
	Withdraw(ctx context.Context, userId int, withdrawCount float64) (*entity.UserDB, error)
}

type handler struct {
	logger      logging2.Logger
	service     Service
	jwtClient   JWTClient
	userExister UserExister
}

type WithdrawRequestBody struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func NewHandler(logger logging2.Logger, updateService Service, jwtClient JWTClient, userExister UserExister) http2.Handler {
	return &handler{
		logger:      logger,
		service:     updateService,
		jwtClient:   jwtClient,
		userExister: userExister,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Get("/api/user/balance", privateRouter.WithPrivateRouter(http.HandlerFunc(h.userBalance), h.logger, h.jwtClient, h.userExister))
	router.Post("/api/user/balance/withdraw", privateRouter.WithPrivateRouter(http.HandlerFunc(h.userBalanceWithdraw), h.logger, h.jwtClient, h.userExister))
}

// userBalance /api/user/balance
func (h handler) userBalance(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userId := auth.GetUserIdFromContext(ctx)

	userDB, err := h.service.GetById(ctx, *userId)

	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	encodedUserBalance, err := encodeUserBalance(userDB)
	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(encodedUserBalance)
	writer.WriteHeader(http.StatusOK)
}

func encodeUserBalance(userDB *entity.UserDB) ([]byte, error) {
	userBalance := &entity.UserBalanceJSON{
		Current:   userDB.Wallet,
		Withdrawn: userDB.Withdrawn,
	}

	return json.Marshal(userBalance)
}

// userBalance /api/user/balance
// Возможные коды ответа:
// - `200` — успешная обработка запроса;
// - `401` — пользователь не авторизован; - в мидлваре выше
// - `402` — на счету недостаточно средств;
// - `422` — неверный номер заказа;
// - `500` — внутренняя ошибка сервера.
func (h handler) userBalanceWithdraw(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userId := auth.GetUserIdFromContext(ctx)
	withdrawRequestBody, err := decodeWithdrawRequestBody(&request.Body)
	// `500` — внутренняя ошибка сервера.
	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	orderNum, err := strconv.Atoi(withdrawRequestBody.Order)
	// `422` — неверный номер заказа;
	if err != nil || !luhn.Valid(orderNum) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	//fmt.Println(withdrawRequestBody.Sum)
	_, err = h.service.Withdraw(ctx, *userId, withdrawRequestBody.Sum)
	// `402` — на счету недостаточно средств;
	if errors.Is(err, user.ErrWithdrawCountGreaterThanUserBalance) {
		h.logger.Info(err)
		writer.WriteHeader(http.StatusPaymentRequired)
		return
	}
	// `500` — внутренняя ошибка сервера.
	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// `200` — успешная обработка запроса;
	writer.WriteHeader(http.StatusOK)
	return
}

func decodeWithdrawRequestBody(body *io.ReadCloser) (*WithdrawRequestBody, error) {
	var withdrawRequestBody WithdrawRequestBody

	decoder := json.NewDecoder(*body)
	err := decoder.Decode(&withdrawRequestBody)
	if err != nil {
		return nil, err
	}

	return &withdrawRequestBody, nil
}
