package api

import "errors"

var (
	ErrRequestInitiate = errors.New("при формировании запроса произошла ошибка")
	ErrRequestDo       = errors.New("при запросе произошла ошибка ")
)
