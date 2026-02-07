package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Load loads the environment variables from the .env file.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}
