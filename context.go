package Application

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"
)

func newContext[TConfig ConfigType, TDatabase DatabaseType]() *appContext[TConfig, TDatabase] {
	return &appContext[TConfig, TDatabase]{
		parent: context.Background(),
		logger: initLogger(),
	}
}

type Context[TConfig ConfigType, TDatabase DatabaseType] interface {
	context.Context

	Log() *zap.Logger
	Database() TDatabase
	Config() *TConfig
}

type appContext[TConfig ConfigType, TDatabase DatabaseType] struct {
	parent context.Context
	logger *zap.Logger
	config *TConfig
	db     TDatabase
	dbConn *sql.DB
}

func (ctx *appContext[TConfig, TDatabase]) Deadline() (deadline time.Time, ok bool) {
	return ctx.parent.Deadline()
}

func (ctx *appContext[TConfig, TDatabase]) Done() <-chan struct{} {
	return ctx.parent.Done()
}

func (ctx *appContext[TConfig, TDatabase]) Err() error {
	return ctx.parent.Err()
}

func (ctx *appContext[TConfig, TDatabase]) Value(key any) any {
	return ctx.parent.Value(key)
}

func (ctx *appContext[TConfig, TDatabase]) Log() *zap.Logger {
	return ctx.logger
}

func (ctx *appContext[TConfig, TDatabase]) Database() TDatabase {
	return ctx.db
}

func (ctx *appContext[TConfig, TDatabase]) DbConn() *sql.DB {
	return ctx.dbConn
}

func (ctx *appContext[TConfig, TDatabase]) Config() *TConfig {
	return ctx.config
}

func (ctx *appContext[TConfig, TDatabase]) withConfig(cfg *TConfig) *appContext[TConfig, TDatabase] {
	ctx.config = cfg
	return ctx
}

func (ctx *appContext[TConfig, TDatabase]) withDb(db TDatabase) *appContext[TConfig, TDatabase] {
	ctx.db = db
	return ctx
}

func (ctx *appContext[TConfig, TDatabase]) withDbConn(conn *sql.DB) *appContext[TConfig, TDatabase] {
	ctx.dbConn = conn
	return ctx
}
