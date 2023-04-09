package Application

import (
	"github.com/caarlos0/env/v7"
)

func MustParseConfig[TConfig ConfigType]() *TConfig {
	cfg, err := ParseConfig[TConfig]()
	if err != nil {
		panic(err)
	}
	return cfg
}

func ParseConfig[TConfig ConfigType]() (*TConfig, error) {
	cfg := new(TConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (a *app[TConfig, TDatabase]) WithConfig(cfg *TConfig) App[TConfig, TDatabase] {
	a.ctx = a.ctx.withConfig(cfg)
	return a
}
