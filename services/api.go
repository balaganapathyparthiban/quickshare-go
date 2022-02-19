package services

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/balaganapathyparthiban/quickshare-go/db"
	"github.com/balaganapathyparthiban/quickshare-go/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/teris-io/shortid"
)

type FileData struct {
	Type     string
	Progress float64
	Path     string
	Expired  string
}

func FileUpload(c *fiber.Ctx) error {
	query := new(struct {
		Name  string `json:"name" validate:"required"`
		Size  int    `json:"size" validate:"required"`
		Type  string `json:"type" valudate:"required"`
		Async bool   `json:"async"`
	})
	c.QueryParser(query)

	if errors := utilities.Validation(query); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	fd := new(FileData)
	fd.Expired = time.Now().Add(time.Hour * 24).Format("January-2-2006-03:04:05PM")

	hash := sha256.New()
	hash.Write([]byte(time.Now().String()))
	id, _ := shortid.Generate()

	reader := c.Context().RequestBodyStream()
	channel := make(chan int)

	go func(rd *io.Reader) {
		file, err := os.Create(fmt.Sprintf("files/%s.%s", query.Name, query.Type))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		fileSize := 0
		buffer := make([]byte, 0, 512*1024)
		for {
			length, err := io.ReadFull(reader, buffer[:cap(buffer)])
			buffer = buffer[:length]
			if err != nil {
				if err == io.EOF {
					break
				}
			}

			fileSize += length

			_, err = file.Write(buffer)
			if err != nil {
				fmt.Println(err)
				break
			}

			fd.Progress = ((float64(fileSize) / float64(query.Size)) * 100) - 1
		}

		rfile, err := os.Open(fmt.Sprintf("files/%s.%s", query.Name, query.Type))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rfile.Close()

		fd.Progress = 100
		fd.Path = fmt.Sprintf("files/%s.%s", query.Name, query.Type)
		fmt.Println(fd)
		fdMarshal, _ := json.Marshal(&fd)

		err = db.Store.Put([]byte(id), fdMarshal, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		channel <- 1
	}(&reader)

	if query.Async {
		time.Sleep(time.Duration(1) * time.Second)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"progress_url": fmt.Sprintf("/api/file/upload/progress?id=%s", id),
		})
	} else {
		<-channel

		return c.JSON(fiber.Map{
			"cid": "",
		})
	}
}

func FileUploadProgress(c *fiber.Ctx) error {
	query := new(struct {
		Id string `json:"id" validate:"required,min=1"`
	})
	c.QueryParser(query)

	if errors := utilities.Validation(query); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	data, err := db.Store.Get([]byte(query.Id), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var fdUnmarshal FileData
	json.Unmarshal(data, &fdUnmarshal)

	return c.JSON(fiber.Map{
		"progress":  fdUnmarshal.Progress,
		"share_url": fmt.Sprintf("/api/file/download?id=%s", query.Id),
		"expired":   fdUnmarshal.Expired,
	})
}

func FileDownload(c *fiber.Ctx) error {
	query := new(struct {
		Id string `json:"id" validate:"required,min=1"`
	})
	c.QueryParser(query)

	if errors := utilities.Validation(query); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	data, err := db.Store.Get([]byte(query.Id), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var fdUnmarshal FileData
	json.Unmarshal(data, &fdUnmarshal)

	f, _ := os.Open(fdUnmarshal.Path)

	return c.Status(fiber.StatusOK).SendStream(bufio.NewReader(f))
}

func ShortenerUrl(c *fiber.Ctx) error {
	query := new(struct {
		Url string `json:"url" validate:"required,min=1"`
	})
	c.QueryParser(query)

	if errors := utilities.Validation(query); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	id, _ := shortid.Generate()

	err := db.Store.Put([]byte(id), []byte(query.Url), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"short_url": fmt.Sprintf("/%s", id),
	})
}

func RedirectUrl(c *fiber.Ctx) error {
	url, error := db.Store.Get([]byte(c.Params("URL")), nil)
	if error != nil {
		return error
	}

	c.Set("location", string(url))
	return c.Redirect(string(url))
}
