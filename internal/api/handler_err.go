package api

import "github.com/gofiber/fiber/v2"

func HandlerErr(f *fiber.Ctx) error {
	return respondWithError(f, 400, "Something went wrong")
}
