package app

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"
	"gopkg.in/reform.v1"
)

func newContext() Context {
	return &appContext{
		parent: context.Background(),
		logger: initLogger(),
	}
}

type Context interface {
	context.Context

	Log() *zap.Logger
	Db() *reform.DB
	DbConn() *sql.DB

	withDb(db *reform.DB) Context
	withDbConn(conn *sql.DB) Context
}

type appContext struct {
	parent context.Context
	logger *zap.Logger
	db     *reform.DB
	dbConn *sql.DB
}

func (ctx *appContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.parent.Deadline()
}

func (ctx *appContext) Done() <-chan struct{} {
	return ctx.parent.Done()
}

func (ctx *appContext) Err() error {
	return ctx.parent.Err()
}

func (ctx *appContext) Value(key any) any {
	return ctx.parent.Value(key)
}

func (ctx *appContext) Log() *zap.Logger {
	return ctx.logger
}

func (ctx *appContext) Db() *reform.DB {
	return ctx.db
}

func (ctx *appContext) DbConn() *sql.DB {
	return ctx.dbConn
}

func (ctx *appContext) withDb(db *reform.DB) Context {
	ctx.db = db
	return ctx
}

func (ctx *appContext) withDbConn(conn *sql.DB) Context {
	ctx.dbConn = conn
	return ctx
}
