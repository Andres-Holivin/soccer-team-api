package utils

import (
	"log"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

type Env struct {
	GinMode              string `env:"GIN_MODE"`
	AppEnv               string `env:"APP_ENV"`
	AppPort              string `env:"APP_PORT"`
	AllowOrigins         string `env:"ALLOW_ORIGINS"`
	DatabaseURL          string `env:"DATABASE_URL"`
	JWTSecret            string `env:"JWT_SECRET"`
	JWTExpiresIn         string `env:"JWT_EXPIRES_IN"`
	CloudinaryURL        string `env:"CLOUDINARY_URL"`
	CloudinaryRootFolder string `env:"CLOUDINARY_ROOT_FOLDER"`
}

func LoadEnv() *Env {
	_ = godotenv.Load()
	var cfg Env
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error loading env: %v", err)
	}
	return &cfg
}
