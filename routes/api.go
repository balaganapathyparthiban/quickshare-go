package routes

import (
	"github.com/balaganapathyparthiban/quickshare-go/services"
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(app *fiber.App) {
	api := app.Group("/api")

	/* File api */
	api.Post("/file/upload", services.FileUpload)
	api.Get("/file/info", services.FileInfo)
	api.Get("/file/download", services.FileDownload)
}
