package config

type Builder interface {
	loadEnvConfig()
	loadFlagsConfig()
	getConfig() Config
}
