package config

type Config struct {
	DBUrl      string
	AuthUrl    string
	PublicKey  string
}

func Load() *Config {
	// TODO: Загрузка из env/config.yaml
	return &Config{}
} 