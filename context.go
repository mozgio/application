package Application

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
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
	Nats() *nats.Conn
	Database() TDatabase
	Config() *TConfig
}

type appContext[TConfig ConfigType, TDatabase DatabaseType] struct {
	parent    context.Context
	logger    *zap.Logger
	nats      *nats.Conn
	config    *TConfig
	db        TDatabase
	closeFunc []func() error
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

func (ctx *appContext[TConfig, TDatabase]) Config() *TConfig {
	return ctx.config
}

func (ctx *appContext[TConfig, TDatabase]) Nats() *nats.Conn {
	return ctx.nats
}

func (ctx *appContext[TConfig, TDatabase]) withConfig(cfg *TConfig) *appContext[TConfig, TDatabase] {
	ctx.config = cfg
	return ctx
}

func (ctx *appContext[TConfig, TDatabase]) withDatabase(db TDatabase) *appContext[TConfig, TDatabase] {
	ctx.db = db
	return ctx
}

func (ctx *appContext[TConfig, TDatabase]) withNats(nc *nats.Conn) *appContext[TConfig, TDatabase] {
	ctx.nats = nc
	return ctx
}
