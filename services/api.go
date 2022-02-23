package services

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
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

func FileUpload(c *fiber.Ctx) error {
	query := new(utilities.FileUpload)
	c.QueryParser(query)

	if errors := utilities.Validation(query); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	fd := new(utilities.FileData)
	fd.Title = query.Title
	fd.Message = query.Message
	if len(query.Password) > 0 {
		passwordHash := sha256.New()
		passwordHash.Write([]byte(query.Password))
		fd.Password = base64.URLEncoding.EncodeToString(passwordHash.Sum(nil))
	}
	fd.Expired = time.Now().Add(time.Hour * 24)

	id, _ := shortid.Generate()

	reader := c.Context().RequestBodyStream()
	channel := make(chan int)

	go func(rd *io.Reader) {
		folderPath := fmt.Sprintf("files/%s", id)
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			os.Mkdir(folderPath, 0777)
		}

		filePath := fmt.Sprintf("files/%s/%s", id, query.Name)
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		buffer := make([]byte, 0, 1024*1024)
		for {
			length, err := io.ReadFull(reader, buffer[:cap(buffer)])
			buffer = buffer[:length]
			if err != nil {
				if err == io.EOF {
					break
				}
			}

			_, err = file.Write(buffer)
			if err != nil {
				fmt.Println(err)
				break
			}
		}

		rfile, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rfile.Close()

		fd.Path = filePath
		fdMarshal, _ := json.Marshal(&fd)

		err = db.Store.Put([]byte(id), fdMarshal, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		channel <- 1
	}(&reader)

	<-channel

	return c.JSON(fiber.Map{
		"file_id": id,
		"expired": fd.Expired,
	})
}

func FileInfo(c *fiber.Ctx) error {
	query := new(utilities.FileDownload)
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

	var fdUnmarshal utilities.FileData
	json.Unmarshal(data, &fdUnmarshal)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"title":               fdUnmarshal.Title,
		"message":             fdUnmarshal.Message,
		"expired":             fdUnmarshal.Expired,
		"path":                fdUnmarshal.Path,
		"isPasswordProtected": len(fdUnmarshal.Password) > 0,
	})
}

func FileDownload(c *fiber.Ctx) error {
	query := new(utilities.FileDownload)
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

	var fdUnmarshal utilities.FileData
	json.Unmarshal(data, &fdUnmarshal)

	fmt.Println(fdUnmarshal)

	if len(fdUnmarshal.Password) > 0 {
		passwordHash := sha256.New()
		passwordHash.Write([]byte(query.Password))
		password := base64.URLEncoding.EncodeToString(passwordHash.Sum(nil))

		if fdUnmarshal.Password != password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid password",
			})
		}
	}

	f, _ := os.Open(fdUnmarshal.Path)

	return c.Status(fiber.StatusOK).SendStream(bufio.NewReader(f))
}
