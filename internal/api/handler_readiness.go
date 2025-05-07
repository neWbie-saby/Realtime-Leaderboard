package api

import "github.com/gofiber/fiber/v2"

func HandlerReadiness(f *fiber.Ctx) error {
	return RespondWithJSON(f, 200, struct{}{})
}
