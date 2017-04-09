package controllers

import (
	"bytes"
	"encoding/base64"
	"github.com/astaxie/beego/logs"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"time"
	"math/rand"
)

var MAX_THUMB uint = 160
var MAX_SHOW uint = 1080

func resizeImage(localPath string, fileName string, format string) (thumb string, goodToShow bool, resizedFileName string, err error) {
	file, err := os.Open(localPath+fileName)
	if err != nil {
		return "", false, "", err
	}

	var img image.Image
	if format == "jpg" {
		img, err = jpeg.Decode(file)
	} else if format == "png" {
		img, err = png.Decode(file)
	} else if format == "gif" {
		img, err = gif.Decode(file)
	}

	if err != nil {
		logs.Warn("Resize image fail")
		return "", false, "", err
	}

	horizontal := img.Bounds().Size().X > img.Bounds().Size().Y
	goodToShow = true
	if img.Bounds().Size().X > int(MAX_SHOW) || img.Bounds().Size().Y > int(MAX_SHOW) {
		goodToShow = false
	}

	var thumbImage image.Image
	if horizontal {
		thumbImage = resize.Resize(MAX_THUMB, 0, img, resize.Lanczos3)
	} else {
		thumbImage = resize.Resize(0, MAX_THUMB, img, resize.Lanczos3)
	}

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, thumbImage, nil)
	thumb = base64.StdEncoding.EncodeToString(buf.Bytes())

	if goodToShow {
		return thumb, goodToShow, "", nil
	}

	var resizeImage image.Image
	if horizontal {
		resizeImage = resize.Resize(MAX_SHOW, 0, img, resize.Lanczos3)
	} else {
		resizeImage = resize.Resize(0, MAX_SHOW, img, resize.Lanczos3)
	}

	resizedFileName =strconv.FormatInt(time.Now().UnixNano() / 1000000,10)+strconv.Itoa(rand.Intn(100)) + ".jpg"
	out, err := os.Create(localPath+resizedFileName)
	if err != nil {
		return "", false, "", err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, resizeImage, nil)
	return thumb, goodToShow, resizedFileName, nil
}
