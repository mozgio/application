package Application

func (a *app[TConfig, TDatabase]) WithSwagger(contents []byte) App[TConfig, TDatabase] {
	a.withSwagger = true
	a.swaggerConfig = swaggerConfig{contents}
	return a
}

type swaggerConfig struct {
	fileContext []byte
}
