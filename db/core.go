package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Code-Hex/sqlx-transactionmanager"
)

func OpenDB(driver, dsn string) *sqlx.DB {
	db := sqlx.MustOpen(driver, dsn)

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(30)

	return db;
}

func (c *Core) Ping(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(*c.Ctx, timeout)
	defer cancel()

	err = c.DB.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("ping db failed: %v", err)
	}

	return
}

func (c *Core) Do(opts *JobOptions, job Job, fail Fail) {
	ctx, cancel := context.WithTimeout(*c.Ctx, opts.Timeout)
	defer cancel()

	if opts.TxOpts == nil {
		opts.TxOpts = &sql.TxOptions{}
	}
	txm, err := c.DB.BeginTxmx(ctx, opts.TxOpts)
	if err != nil {
		fail(err)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			txm.Rollback()
			fail(fmt.Errorf("%v", err))
		}
	}()

	job(&Core{DB: c.DB, Txm: txm, Ctx: &ctx})

	txm.Commit()
}

func (c *Core) DoSimple(opts *JobOptions, job Job, fail Fail) {
	ctx, cancel := context.WithTimeout(*c.Ctx, opts.Timeout)
	defer cancel()

	defer func() {
		if err := recover(); err != nil {
			fail(fmt.Errorf("%v", err))
		}
	}()

	job(&Core{DB: c.DB, Ctx: &ctx})
}

func NewJobOpts(simple, auto bool) *JobOptions {
	return &JobOptions{ Timeout: 5 * time.Second, Simple: simple, Auto: auto }
}
