package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServerConfig struct {
	Addr            string `yaml:"address"`
	ReadTimeout     string `yaml:"read_timeout"`
	WriteTimeout    string `yaml:"write_timeout"`
	IdleTimeout     string `yaml:"idle_timeout"`
	ShutdownTimeout string `yaml:"shutdown_timeout"`
}

type Config struct {
	Env         string           `yaml:"env" env:"Env" env-required:"true"`
	StoragePath string           `yaml:"storage_path" env-required:"true"`
	HTTP        HttpServerConfig `yaml:"http_server"`
}

func MustLoadConfig() *Config {
	var configpath string

	configpath = os.Getenv("CONFIG_PATH")

	if configpath == "" {
		flags := flag.String("config", "", "Path to config file")
		flag.Parse()

		configpath = *flags

		if configpath == "" {
			log.Fatal("CONFIG_PATH environment variable or --config flag must be set")
		}
	}

	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configpath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configpath, &cfg)

	if err != nil {
		log.Fatalf("Failed to load config: %v", err.Error())
	}

	return &cfg

}
