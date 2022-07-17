package config

type Config struct {
	LogLevel string `toml:"log_level"`
	LogAddr  string `toml:"log_path"`
	Domain   string `toml:"domain"`
	BindAddr string `toml:"bind_addr"`
	MongoURL string `toml:"mongo_url"`
}
