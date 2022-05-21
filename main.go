package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	validator "github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port   int    `default:"8000"`
	Host   string `default:"localhost"`
	Scheme string `default:"https"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()

	app.Use(cache.New(cache.Config{
		Expiration:   30 * 24 * time.Hour,
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.OriginalURL())
		},
	}))

	app.Use(limiter.New(limiter.Config{
		Max:               5,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Jakarta",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("👋")
	})

	app.Get("/compress/*", func(c *fiber.Ctx) error {
		requestedUrl := strings.TrimPrefix(c.OriginalURL(), "/compress/")
		requestedUrl = strings.TrimPrefix(requestedUrl, "/compress?url=")

		isValid := validator.IsURL(requestedUrl)
		if !isValid {
			return c.Status(fiber.StatusBadRequest).SendString("Not a valid URL")
		}

		generatedUrl := url.URL{
			Scheme: config.Scheme,
			Host:   config.Host,
			Path:   compress(requestedUrl),
		}
		return c.SendString(generatedUrl.String())
	})

	app.Get("/u/*", func(c *fiber.Ctx) error {
		return c.Redirect(decompress(c.Params("*")), fiber.StatusMovedPermanently)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.Port)))
}
