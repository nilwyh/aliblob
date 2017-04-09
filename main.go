package main

import (
	"github.com/astaxie/beego"
	_ "github.com/nilwyh/aliblob/routers"
)

func main() {
	beego.Run()
}
