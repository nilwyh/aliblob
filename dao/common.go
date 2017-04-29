package dao

import (
	"gopkg.in/olivere/elastic.v5"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

var EsClient *elastic.Client = nil
var EsMediaIndex = beego.AppConfig.String("esindexmedia")

func InitEsClient() {
	var err error
	EsClient, err = elastic.NewSimpleClient(elastic.SetURL(beego.AppConfig.String("metadataserver")))
	if err != nil {
		logs.Warn(err)
		// Handle error
		panic(err)
	}
}