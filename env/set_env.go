package env

import (
	"github.com/joho/godotenv"
)

func Load_env(path string) {
	err := godotenv.Load(path)

	if err != nil {
		panic("Error loading specified env file")
	}
}
