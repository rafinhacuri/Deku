package env

import (
	"os"
)

type EnvConfig struct {
	JWT_SECRET string
	CRYPT_KEY  string
}

var C *EnvConfig

func LoadEnv() {
	C = &EnvConfig{
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		CRYPT_KEY:  os.Getenv("CRYPT_KEY"),
	}
}
