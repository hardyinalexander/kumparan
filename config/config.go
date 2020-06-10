package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port       string `envconfig:"PORT" default:"8080"`
	DBHost     string `envconfig:"DB_HOST"`
	DBPort     string `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER"`
	DBPassword string `envconfig:"DB_HOST"`
	DBName     string `envconfig:"DB_NAME"`
}

// Get Configuration function
func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return cfg
}
