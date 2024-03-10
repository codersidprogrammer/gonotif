package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ExitIfErr(err error, message ...interface{}) error {
	if err != nil {
		log.Fatal(message, err)
	}
	return nil
}

func MessageIfErr(err error, message ...interface{}) error {
	if err != nil {
		log.Fatal(message, err)
		return err
	}

	return err
}

func ReturnErrorIfErr(d any, err error) (any, error) {
	if err != nil {
		return nil, err
	}

	return d, nil
}

func ReturnErrMessageIfErr(err error, message interface{}, ctx *fiber.Ctx) error {
	if err != nil {
		return ctx.Status(fiber.ErrBadGateway.Code).JSON(&fiber.Map{
			"message": message,
			"error":   err,
		})
	}

	return nil
}

func ReturnHttpErr400MessageIfErr(err error, message interface{}, ctx *fiber.Ctx) error {
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(&fiber.Map{
			"message": message,
			"error":   err,
		})
	}

	return nil
}

func ReturnHttpErr404MessageIfErr(err error, message interface{}, ctx *fiber.Ctx) error {
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(&fiber.Map{
			"message": message,
			"error":   err,
		})
	}

	return nil
}
