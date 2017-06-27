package controllers

import (
	"github.com/astaxie/beego"
	"os"
	"time"
	"strings"
	"github.com/nilwyh/aliblob/dao"
	"mime/multipart"
	"github.com/astaxie/beego/logs"
)

var ERROR_REDIRECT string = beego.AppConfig.String("errorredirect")
var SUCCESS_REDIRECT string = beego.AppConfig.String("successredirect")

var IMAGE map[string]string = map[string]string{".png":"png",".jpg":"jpg",".bmp":"bmp",".jpeg":"jpeg",".gif":"gif"}
var VIDEO map[string]string = map[string]string{".avi":"avi",".flv":"flv",".wmv":"wmv",".mov":"mov",".mp4":"mp4",".mkv":"mkv"}

type UploadController struct {
	beego.Controller
	localpath string
	urlpath string
}

func (this *UploadController) Post() {
	files, err := this.GetFiles("myfiles")
	//TODO: trim /
	collection := "default"
	if mc := this.GetString("mycollect"); mc !="" {
		collection = this.GetString("mycollect")
	}
	authGroup:= this.GetString("myauthgroup")
	//year,month,_:= time.Now().Date()
	//shard:=strconv.Itoa(year)+"_"+month.String()+"/"
	root:=beego.AppConfig.String("blobroot")
	this.localpath=root + collection + "/"
	this.urlpath=collection + "/"
	if err != nil {
		this.Redirect(ERROR_REDIRECT, 302)
		return
	}
	if len(files) > 0 {
		os.MkdirAll(this.localpath, 0777)
	}
	for i, _ := range files {
		name := "file_" + files[i].Filename
		timestamp,err1:=this.GetInt64(name)
		if err1!=nil {
			timestamp = time.Now().UnixNano()/1000000
		}
		modifiedTime := time.Unix(0, timestamp * int64(time.Millisecond))
		this.handleFile(files[i], modifiedTime, collection, authGroup)
	}
	this.Redirect(SUCCESS_REDIRECT, 302)
}

func (this *UploadController) handleFile(f *multipart.FileHeader, modifiedTime time.Time, collect, authGroup string) {
	suffix := getSuffix(f.Filename)
	if _, ok := IMAGE[suffix]; ok {
		this.handleImage(f, suffix, modifiedTime, collect, authGroup)
	} else if _, ok := VIDEO[suffix]; ok {
		this.handleVideo(f, suffix, modifiedTime)
	} else {
		this.handleOther(f, suffix)
	}
}

func (this *UploadController) handleImage(f *multipart.FileHeader, suffix string, modifiedTime time.Time, collect, authGroup string) {
	//for each fileheader, get a handle to the actual file
	fileName, hashString, err:= dao.CopyFileMultipart(f, this.localpath, suffix)
	if err!=nil {
		this.Redirect(ERROR_REDIRECT, 302)
		logs.Warn("err:%s",err)
		return
	}

	thumb, goodToShow, resizedFileName, err := resizeImage(this.localpath, fileName, IMAGE[suffix])
	if err!=nil {
		this.Redirect(ERROR_REDIRECT, 302)
		logs.Warn("err:%s",err)
		return
	}

	rawId, err:=dao.UpdateImageMetadata(this.urlpath+fileName, hashString, IMAGE[suffix], thumb, goodToShow, true, "", modifiedTime, collect, authGroup)
	if err!=nil {
		this.Redirect(ERROR_REDIRECT, 302)
		logs.Warn("err:%s",err)
		return
	}
	if !goodToShow {
		dao.UpdateImageMetadata(this.urlpath+resizedFileName, hashString, IMAGE[suffix], thumb, true, false, rawId, modifiedTime, collect, authGroup)
	}
}

func (this *UploadController) handleVideo(f *multipart.FileHeader, suffix string, modifiedTime time.Time) {
	fileName, hashString, err:= dao.CopyFileMultipart(f, this.localpath, suffix)
	if err!=nil {
		this.Redirect(ERROR_REDIRECT, 302)
		logs.Warn("err:%s",err)
		return
	}

	dao.UpdateVideoMetadata(this.urlpath+fileName, hashString, VIDEO[suffix], modifiedTime)
}

func (this *UploadController) handleOther(f *multipart.FileHeader, suffix string) {
	fileName, hashString, err:= dao.CopyFileMultipart(f, this.localpath, suffix)
	if err!=nil {
		this.Redirect(ERROR_REDIRECT, 302)
		logs.Warn("err:%s",err)
		return
	}

	dao.UpdateOtherMetadata(this.urlpath+fileName, hashString, suffix)
}

func getSuffix(name string) (string) {
	index:=strings.LastIndex(name, ".")
	if index ==-1 {
		return ""
	}
	return strings.ToLower(name[index:])
}
