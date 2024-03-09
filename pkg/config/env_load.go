package config

import (
	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/joho/godotenv"
)

func LoadEnvironment(fp *string) error {
	err := godotenv.Load(*fp)
	utils.ExitIfErr(err, "Failed to load environment")
	return nil
}
