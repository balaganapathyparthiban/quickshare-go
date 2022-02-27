package main

import (
	"flag"
	"log"
	"os"

	"github.com/balaganapathyparthiban/quickshare-go/db"
	"github.com/balaganapathyparthiban/quickshare-go/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"
	"github.com/qinains/fastergoding"
)

var IsDev bool

func init() {
	/* Check --dev is passed or not */
	flag.BoolVar(&IsDev, "dev", false, "Pass --dev to load .env.dev file")
	flag.Parse()

	/* In development get env values */
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatalf("Error loading .env.dev file")
	}

	if IsDev {
		/* Enable hotreload */
		fastergoding.Run()
	}

	if _, err := os.Stat("files"); os.IsNotExist(err) {
		os.Mkdir("files", 0777)
	}

	db.InitDB()
}

func main() {
	app := fiber.New()

	app.Server().StreamRequestBody = true

	/* Middlewares */
	app.Use(cors.New())
	app.Use(helmet.New())

	/* Routes */
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Quick Share")
	})

	routes.ApiRoutes(app)

	app.All("*", func(c *fiber.Ctx) error {
		return c.SendString("Not an valid path/method.")
	})

	app.Listen(os.Getenv("PORT"))
}
