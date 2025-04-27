package api

import "github.com/gofiber/fiber/v2"

func respondWithJSON(f *fiber.Ctx, code int, payload interface{}) error {
	f.Status(code)
	return f.JSON(payload)
}
