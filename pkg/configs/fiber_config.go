package configs

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

// FiberConfig returns a new fiber Config
// for the API.
func FiberConfig() fiber.Config {
	return fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		ReadTimeout: 30 * time.Second,
	}
}
