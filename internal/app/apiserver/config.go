package apiserver

type Config struct {
	BindAddr       string `json:"bind_addr"`
	LogLevel       string `json:"log_level"`
	DatabaseURL    string `json:"database_url"`
	RedisLocalhost string `json:"redislocalhost"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:       ":8080",
		LogLevel:       "debug",
		RedisLocalhost: "6379",
	}
}
