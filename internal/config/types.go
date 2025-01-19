package config

type AppEnvironment int

const (
	Local AppEnvironment = iota
	Development
	Production
)
