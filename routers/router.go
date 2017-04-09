package routers

import (
	"github.com/astaxie/beego"
	"github.com/nilwyh/aliblob/controllers"
)

func init() {
	beego.BConfig.WebConfig.DirectoryIndex=true

	beego.Router("/blob/upload", &controllers.UploadController{})
	beego.SetStaticPath("/blob/get",beego.AppConfig.String("blobroot"))
}
