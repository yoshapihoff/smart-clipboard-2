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

// Функция для декодирования base64 строки в изображение
func DecodeBase64Image(base64String string) (image.Image, error) {
	// Удаляем префикс data:image/png;base64, если он есть
	if strings.HasPrefix(base64String, "data:image") {
		parts := strings.Split(base64String, ",")
		if len(parts) > 1 {
			base64String = parts[1]
		}
	}

	// Декодируем base64
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	// Декодируем PNG изображение
	img, err := png.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Функция для кодирования изображения в байты
func EncodeImageToBytes(img image.Image) []byte {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

// Функция для создания ресурса изображения из base64
func CreateImageResourceFromBase64(base64String string) *fyne.StaticResource {
	img, err := DecodeBase64Image(base64String)
	if err != nil {
		fmt.Printf("Ошибка декодирования иконки: %v\n", err)
		return nil
	}

	return fyne.NewStaticResource("icon", EncodeImageToBytes(img))
}
