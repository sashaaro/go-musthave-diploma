package sql

import (
	"context"
	"errors"
	"github.com/GTech1256/go-musthave-diploma-tpl/internal/domain/entity"
	"github.com/GTech1256/go-musthave-diploma-tpl/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type storage struct {
	logger logging.Logger
	db     DB
}

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Ping(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

var (
	ErrNotUniqueLogin                      = errors.New("пользователь с таким логином уже зарегистрирован")
	ErrInvalidLoginPasswordCombination     = errors.New("неверная пара логин/пароль")
	ErrWithdrawCountGreaterThanUserBalance = errors.New("Запрошенная сумма вывода больше, чем баланс пользователя")
)

func NewStorage(db DB, logger logging.Logger) *storage {
	return &storage{
		db:     db,
		logger: logger,
	}
}

func (s *storage) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *storage) Register(ctx context.Context, userRegister *entity.UserRegisterJSON) (*entity.UserDB, error) {
	var userDB entity.UserDB
	err := s.db.QueryRow(
		ctx,
		"INSERT INTO gophermart.users (login, password) values ($1, $2) RETURNING id, login, password",
		userRegister.Login,
		userRegister.Password,
	).Scan(&userDB.ID, &userDB.Login, &userDB.Password)

	// Проверка на то что логин не уникальный
	if err != nil {
		if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
			return nil, ErrNotUniqueLogin
		}
	}

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &userDB, nil
}

func (s *storage) Login(ctx context.Context, userRegister *entity.UserLoginJSON) (*entity.UserDB, error) {
	var userDB entity.UserDB

	err := s.db.QueryRow(
		ctx,
		"SELECT id, login, password FROM gophermart.users WHERE login = $1 AND password = $2",
		userRegister.Login,
		userRegister.Password,
	).Scan(&userDB.ID, &userDB.Login, &userDB.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrInvalidLoginPasswordCombination
	}
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &userDB, nil
}

func (s *storage) GetById(ctx context.Context, userId int) (*entity.UserDB, error) {
	var userDB entity.UserDB

	err := s.db.QueryRow(
		ctx,
		"SELECT id, login, password, wallet, withdrawn FROM gophermart.users WHERE id = $1",
		userId,
	).Scan(&userDB.ID, &userDB.Login, &userDB.Password, &userDB.Wallet, &userDB.Withdrawn)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &userDB, nil
}

func (s *storage) Withdraw(ctx context.Context, userId int, withdrawCount float64) (*entity.UserDB, error) {
	var userDB entity.UserDB

	err := s.db.QueryRow(
		ctx,
		"SELECT id, login, password, wallet, withdrawn FROM gophermart.users WHERE id = $1",
		userId,
	).Scan(&userDB.ID, &userDB.Login, &userDB.Password, &userDB.Wallet, &userDB.Withdrawn)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if userDB.Wallet < withdrawCount {
		return nil, ErrWithdrawCountGreaterThanUserBalance
	}

	err = s.db.QueryRow(
		ctx,
		"UPDATE gophermart.users SET wallet=$2, withdrawn=$3 WHERE id = $1 RETURNING wallet, withdrawn",
		userId,
		userDB.Wallet-withdrawCount,
		userDB.Withdrawn+withdrawCount,
	).Scan(&userDB.Wallet, &userDB.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &userDB, nil
}

func (s *storage) IncrementBalance(ctx context.Context, userId int, incValue float64) (*entity.UserDB, error) {
	var userDB entity.UserDB

	err := s.db.QueryRow(
		ctx,
		"UPDATE gophermart.users SET wallet = wallet + $2 WHERE id = $1 RETURNING id, login, password, wallet, withdrawn",
		userId,
		incValue,
	).Scan(&userDB.ID, &userDB.Login, &userDB.Password, &userDB.Wallet, &userDB.Withdrawn)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &userDB, nil
}
