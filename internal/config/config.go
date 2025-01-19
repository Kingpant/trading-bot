package config

import (
	"log"
	"reflect"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	AppEnv AppEnvironment `env:"APP_ENV" envDefault:"local"`

	BitkubBaseURL   string `env:"BITKUB_BASE_URL" envDefault:"https://api.bitkub.com"`
	BitkubAPIKey    string `env:"BITKUB_API_KEY,required"`
	BitkubAPISecret string `env:"BITKUB_API_SECRET,required"`
}

func LoadConfig(envPath string) (*Config, error) {
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			return nil, errors.Wrap(err, "failed to load .env file")
		}
	} else {
		if err := godotenv.Load(); err != nil {
			return nil, errors.Wrap(err, "failed to load .env file")
		}
	}

	var cfg Config
	if err := env.ParseWithOptions(&cfg, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(AppEnvironment(0)): validateAppEnvironment,
		},
	}); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	log.Println("Config loaded successfully")
	return &cfg, nil
}
