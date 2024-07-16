package entity

import (
	"database/sql"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrderNumber int

type OrderDB struct {
	ID         int
	Number     int
	Status     string
	Accrual    sql.NullFloat64
	UploadedAt pgtype.Timestamp
	UserId     int
}
