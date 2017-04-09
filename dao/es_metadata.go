package dao

import (
	"context"
	"github.com/nilwyh/aliblob/models"
	"time"
)


func UpdateImageMetadata(src, sha256string, format string, thumb string, goodToShow bool, isRaw bool, rawId string, modifiedTime time.Time) (string, error) {
	if EsClient == nil {
		InitEsClient()
	}
	ctx := context.Background()

	image:=models.Image{
		Src:src,
		CreatedTime: modifiedTime,
		SHA256: sha256string,
		Thumb: thumb,
		GoodToShow:goodToShow,
		Format:format,
		IsRaw: isRaw,
		RawId:rawId,
	}
	resp, err := EsClient.Index().
		Index("media").
		Type("image").
		BodyJson(image).
		Refresh("true").
		Do(ctx)
	if err!=nil {
		return "", err
	}
	return resp.Id, nil
}

func UpdateVideoMetadata(src string, sha256string string, format string, modifiedTime time.Time) (string, error) {
	if EsClient == nil {
		InitEsClient()
	}
	ctx := context.Background()

	video:=models.Video{
		Src:src,
		CreatedTime:modifiedTime,
		SHA256: sha256string,
		Format: format,
	}

	resp, err := EsClient.Index().
		Index("media").
		Type("video").
		BodyJson(video).
		Refresh("true").
		Do(ctx)
	if err!=nil {
		return "", err
	}
	return resp.Id, nil
}

func UpdateOtherMetadata(src string, sha256string string, suffix string) (string, error) {
	format := ""
	if suffix != "" {
		format = suffix[1:]
	}

	if EsClient == nil {
		InitEsClient()
	}
	ctx := context.Background()

	other:=models.Other{
		Src:src,
		CreatedTime:time.Now(),
		SHA256: sha256string,
		Format: format,
	}

	resp, err := EsClient.Index().
		Index("media").
		Type("other").
		BodyJson(other).
		Refresh("true").
		Do(ctx)
	if err!=nil {
		return "", err
	}
	return resp.Id, nil
}
