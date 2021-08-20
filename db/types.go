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
	Job     Job
	Fail    Fail
}
