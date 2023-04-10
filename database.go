package application

import (
	"io/fs"

	"github.com/mozgio/database"
)

func (a *app[TConfig, TDatabase]) WithDatabase(driver database.Driver[TDatabase]) App[TConfig, TDatabase] {
	a.withDatabase = true
	a.databaseDriver = driver
	return a
}

func (a *app[TConfig, TDatabase]) WithMigrations(migrationsFs fs.FS, pattern string) App[TConfig, TDatabase] {
	a.withMigrations = true
	a.migrationsConfig = migrationsConfig{migrationsFs, pattern}
	return a
}

type migrationsConfig struct {
	fs      fs.FS
	pattern string
}
