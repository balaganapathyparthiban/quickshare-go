package routes

import (
	"github.com/balaganapathyparthiban/quickshare-go/services"
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(app *fiber.App) {
	api := app.Group("/api")

	/* File api */
	api.Post("/file/upload", services.FileUpload)
	api.Get("/file/upload/progress", services.FileUploadProgress)
	api.Get("/file/download", services.FileDownload)

	/* Shortener api */
	api.Get("/shorten/url", services.ShortenUrl)
}
