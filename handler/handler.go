package handler

import (
	"semprojdb/db/psql"
	"semprojdb/logger"
)

type Handler struct {
	db *psql.DB
	lg *logger.Logger
}

func NewHandler(db *psql.DB, lg *logger.Logger) *Handler {
	return &Handler{
		db: db,
		lg: lg,
	}
}

type WithError struct {
	Err string      `json:"error"`
	Val interface{} `json:"value"`
}
