package utilities

import (
	"time"
)

type FileData struct {
	Title    string
	Message  string
	Password string
	Name     string
	Path     string
	Size     int
	Expired  time.Time
}

type FileUpload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Password string `json:"password"`
	Name     string `json:"name" validate:"required"`
	Size     int    `json:"size" validate:"required,number,gte=1,lte=1073741824"`
}

type FileDownload struct {
	Id       string `json:"id" validate:"required,min=1"`
	Password string `json:"password"`
}
