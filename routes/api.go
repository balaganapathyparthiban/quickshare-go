package routes

import (
	"github.com/balaganapathyparthiban/quickshare-go/services"
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(router fiber.Router) {
	/* File api */
	router.Post("/file/upload", services.FileUpload)
	router.Get("/file/upload/progress", services.FileUploadProgress)
	router.Get("/file/download", services.FileDownload)

	/* Shortener api */
	router.Get("/shortener/url", services.ShortenerUrl)
}
