package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/Code-Hex/sqlx-transactionmanager"
)

type Core struct {
	DB  *sqlx.DB
	Txm *sqlx.Txm
	Ctx *context.Context
}

type Job func(*Core)

type Fail func(error)

type JobOptions struct {
	Timeout time.Duration
	TxOpts  *sql.TxOptions
	Simple  bool
	Auto    bool
}

type PageReq struct {
	Target int64 `form:"target" binding:"required"`
	Size   int64 `form:"size" binding:"required"`
}

type PageRes struct {
	TotalRecords int64 `db:"total" json:"total"`
}
