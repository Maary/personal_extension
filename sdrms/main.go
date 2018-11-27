package main

import (
	"github.com/astaxie/beego/plugins/cors"
	_ "personal_extension/sdrms/routers"
	_ "personal_extension/sdrms/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	//orm.RunSyncdb("default", false, true)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.Run()
}
