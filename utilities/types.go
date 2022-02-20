package utilities

import "time"

type FileData struct {
	Progress float64
	Path     string
	Dir      string
	Expired  time.Time
}

type FileUpload struct {
	Name  string `json:"name" validate:"required"`
	Size  int    `json:"size" validate:"required,number,gte=1,lte=1073741824"`
	Type  string `json:"type" validate:"required"`
	Async bool   `json:"async"`
}

type FileUploadProgress struct {
	Id string `json:"id" validate:"required,min=1"`
}

type FileDownload struct {
	Id string `json:"id" validate:"required,min=1"`
}

type ShortenUrl struct {
	Url string `json:"url" validate:"required,min=1"`
}
