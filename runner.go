package application

func (a *app[TConfig, TDatabase]) WithRunner(runnerFunc RunnerFunc[TConfig, TDatabase]) App[TConfig, TDatabase] {
	a.withRunners = true
	a.runners = append(a.runners, runnerFunc)
	return a
}
