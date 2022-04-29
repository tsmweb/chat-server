package imageutil

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"mime/multipart"
)

func CreateImage() *image.RGBA {
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

func CreateImageBuffer(imageType string) (contentType string, content *bytes.Buffer, err error) {
	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		return "", nil, err
	}

	img := CreateImage()

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
