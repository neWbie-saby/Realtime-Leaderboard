package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger(f *fiber.Ctx) error {
	start := time.Now()
	err := f.Next()
	stop := time.Now()

	log.Printf("[%s] %s %s (%v)", f.Method(), f.Path(), f.IP(), stop.Sub(start))
	return err
}
