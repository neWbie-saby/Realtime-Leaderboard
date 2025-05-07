package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func RespondWithJSON(f *fiber.Ctx, code int, payload interface{}) error {
	f.Status(code)
	return f.JSON(payload)
}

func RespondWithError(f *fiber.Ctx, code int, msg string) error {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	return RespondWithJSON(f, code, errResponse{
		Error: msg,
	})
}
