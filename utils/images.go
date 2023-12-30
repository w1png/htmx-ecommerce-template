package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

func CreateImagesFolderIfNotExists() error {
	if _, err := os.Stat("static/images"); os.IsNotExist(err) {
		return os.Mkdir("static/images", 0755)
	}

	return nil
}

func generateFileName(filename string) string {
	splt := strings.Split(filename, ".")
	extension := splt[len(splt)-1]

	for {
		hash_input := filename + fmt.Sprintf("%d", os.Getpid()) + fmt.Sprintf("%d", time.Now().Unix())
		var hash string = fmt.Sprintf("%x", sha256.Sum256([]byte(hash_input)))
		filename := hash[:6] + hash[len(hash)-6:] + "." + extension

		if _, err := os.Stat("static/images/" + filename); os.IsNotExist(err) {
			return filename
		}
	}
}

func SaveImage(img *multipart.FileHeader) (string, error) {
	file, err := img.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	splt := strings.Split(img.Filename, ".")
	extension := splt[len(splt)-1]

	img.Filename = generateFileName("." + extension)
	dst, err := os.Create("static/images/" + img.Filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if _, err = dst.Write(fileBytes); err != nil {
		return "", err
	}

	return img.Filename, nil
}
