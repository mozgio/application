package app

func (a *app) WithSwagger(contents []byte) App {
	a.withSwagger = true
	a.swaggerConfig = swaggerConfig{contents}
	return a
}

type swaggerConfig struct {
	fileContext []byte
}
