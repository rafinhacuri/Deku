package env

import (
	"os"
)

type EnvConfig struct {
	JWT_SECRET  string
	CRYPT_KEY   string
	SQLITE_PATH string
}

var C *EnvConfig

func LoadEnv() {
	C = &EnvConfig{
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
		CRYPT_KEY:   os.Getenv("CRYPT_KEY"),
		SQLITE_PATH: os.Getenv("SQLITE_PATH"),
	}
}
