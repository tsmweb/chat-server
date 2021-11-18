package handler

import (
	"bytes"
	"errors"
	"github.com/tsmweb/file-service/config"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var ErrUnsupportedMediaType = errors.New("unsupported media type")

func validateFileSize(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadSize())
	return r.ParseMultipartForm(config.MaxUploadSize())
}

func getContentType(file multipart.File) (fileType string, fileExtension string, err error) {
	fileHeader := make([]byte, 512)
	if _, err = file.Read(fileHeader); err != nil {
		return
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return
	}

	fileType = http.DetectContentType(fileHeader)

	switch fileType {
	case "image/jpeg", "image/jpg":
		fileExtension = "jpg"

	case "image/png":
		fileExtension = "png"

	case "video/mp4":
		fileExtension = "mp4"

	case "application/pdf":
		fileExtension = "pdf"

	default:
		err = ErrUnsupportedMediaType
	}

	return
}

func copyFile(dstPath string, srcFile io.Reader) error {
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func createImage() *image.RGBA {
	width := 800
	height := 400

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	cyan := color.RGBA{R: 100, G: 200, B: 200, A: 0xff}
	gray := color.RGBA{R: 100, G: 100, B: 100, A: 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2:
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2:
				img.Set(x, y, gray)
			default:
			}
		}
	}

	return img
}

func createImageBuffer(imageType string) (contentType string, content *bytes.Buffer, err error) {
	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		return "", nil, err
	}

	img := createImage()

	if imageType == "png" {
		if err = png.Encode(part, img); err != nil {
			return "", nil, err
		}
	} else {
		if err = jpeg.Encode(part, img, nil); err != nil {
			return "", nil, err
		}
	}

	if err = writer.WriteField("id", "be49afd2ee890805c21ddd55879db1387aec9751"); err != nil {
		return "", nil, err
	}

	contentType = writer.FormDataContentType()
	content = body

	return
}
