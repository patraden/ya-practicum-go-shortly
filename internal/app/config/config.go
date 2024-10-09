package config

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func DefaultConfig() Config {
	return Config{
		ServerAddr: `localhost:8080`,
		BaseURL:    `http://localhost:8080/`,
	}
}

func LoadConfig() Config {
	builder := newDefaultBuilder()
	builder.loadEnvConfig()
	builder.loadFlagsConfig()
	cfg := builder.getConfig()
	return cfg
}
