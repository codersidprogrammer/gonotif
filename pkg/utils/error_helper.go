package utils

import "github.com/gofiber/fiber/v2/log"

func ExitIfErr(err error, message ...interface{}) error {
	if err != nil {
		log.Fatal(message, err)
	}
	return nil
}
