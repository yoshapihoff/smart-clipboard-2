package fyne_app

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"strings"

	"fyne.io/fyne/v2"
)

func DecodeBase64Image(base64String string) (image.Image, error) {
		if strings.HasPrefix(base64String, "data:image") {
		parts := strings.Split(base64String, ",")
		if len(parts) > 1 {
			base64String = parts[1]
		}
	}

		data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

		img, err := png.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func EncodeImageToBytes(img image.Image) []byte {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func CreateImageResourceFromBase64(base64String string) *fyne.StaticResource {
	img, err := DecodeBase64Image(base64String)
	if err != nil {
		fmt.Printf("Ошибка декодирования иконки: %v\n", err)
		return nil
	}

	return fyne.NewStaticResource("icon", EncodeImageToBytes(img))
}
