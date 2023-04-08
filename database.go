package app

import (
	"database/sql"
	"io/fs"
	"time"

	"github.com/adlio/schema"
	"go.uber.org/zap"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/mysql"
)

func (a *app) WithDatabase(dsn string) App {
	a.databaseConfig = databaseConfig{dsn}
	return a
}

func (a *app) WithMigrations(migrationsFs fs.FS, pattern string) App {
	a.migrationsConfig = migrationsConfig{migrationsFs, pattern}
	return a
}

type databaseConfig struct {
	dsn string
}

type migrationsConfig struct {
	fs      fs.FS
	pattern string
}

func connectToDatabase(ctx Context, cfg databaseConfig) Context {
	conn, err := sql.Open("mysql", cfg.dsn)
	if err != nil {
		ctx.Log().Fatal("database connect error", zap.Error(err))
	}
	rf := reform.NewDB(conn, mysql.Dialect, &dbLogger{ctx.Log()})
	return ctx.withDbConn(conn).withDb(rf)
}

func migrateDatabase(ctx Context, cfg migrationsConfig) {
	migrations, err := schema.FSMigrations(cfg.fs, cfg.pattern)
	if err != nil {
		ctx.Log().Fatal("failed to read database migrations", zap.Error(err))
	}
	opts := []schema.Option{
		schema.WithDialect(schema.MySQL),
	}
	if err = schema.NewMigrator(opts...).Apply(ctx.DbConn(), migrations); err != nil {
		ctx.Log().Fatal("failed to migrate database", zap.Error(err))
	}
}

type dbLogger struct {
	log *zap.Logger
}

func (l *dbLogger) Before(query string, args []any) {
	l.log.Debug("before query",
		zap.String("query", query),
		zap.Any("args", args),
	)
}

func (l *dbLogger) After(query string, args []interface{}, d time.Duration, err error) {
	l.log.Debug("after query",
		zap.String("query", query),
		zap.Any("args", args),
		zap.Duration("duration", d),
		zap.Error(err),
	)
}
