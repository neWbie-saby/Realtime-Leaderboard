package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func respondWithJSON(f *fiber.Ctx, code int, payload interface{}) error {
	f.Status(code)
	return f.JSON(payload)
}

func respondWithError(f *fiber.Ctx, code int, msg string) error {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	return respondWithJSON(f, code, errResponse{
		Error: msg,
	})
}
