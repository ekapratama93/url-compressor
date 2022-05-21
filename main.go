package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
    Port int
    Host string
}

func main() {
	var config Config
    err := envconfig.Process("", &config)
	if err != nil {
        log.Fatal(err.Error())
    }

    app := fiber.New()

	app.Use(cache.New(cache.Config{
		Expiration: 30 * 24 * time.Hour,
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.OriginalURL())
		},
	}))

	app.Use(limiter.New(limiter.Config{
		Max:            5,
		Expiration:     30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(logger.New(logger.Config{
		TimeZone:   "Asia/Jakarta",
	}))

	app.Get("/", func (c *fiber.Ctx) error {
		return c.SendString("👋")
	})

    app.Get("/compress/*", func (c *fiber.Ctx) error {
		log.Println(c.OriginalURL())
		url := strings.TrimPrefix(c.OriginalURL(), "/compress/")
		url = strings.TrimPrefix(url, "/compress?url=")
		return c.SendString("https://"+config.Host+"/u/"+compress(url))
	})

	app.Get("/u/*", func (c *fiber.Ctx) error {
		return c.Redirect(decompress(c.Params("*")), fiber.StatusMovedPermanently)
	})

    log.Fatal(app.Listen(fmt.Sprintf(":%d", config.Port)))
}
