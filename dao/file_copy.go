package dao

import (
	"mime/multipart"
	"strconv"
	"os"
	"io"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"math/rand"
)

func CopyFileMultipart(f *multipart.FileHeader, directory string, suffix string) (fileName string, sha256string string, err error) {
	file, err := f.Open()
	defer file.Close()
	if err != nil {
		return "", "", err
	}

	fileName =strconv.FormatInt(time.Now().UnixNano() / 1000000,10)+strconv.Itoa(rand.Intn(100)) + suffix
	filePath := directory+fileName

	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer dst.Close()

	if err != nil {
		return "", "", err
	}
	if _, err := io.Copy(dst, file); err != nil {
		return "", "", err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", "", err
	}
	//Get the 32 bytes hash
	hashInBytes := hasher.Sum(nil)[:32]
	//Convert the bytes to a string
	sha256string = hex.EncodeToString(hashInBytes)

	return fileName, sha256string, nil
}