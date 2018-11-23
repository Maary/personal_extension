package main

import (
	"github.com/astaxie/beego/orm"
	_ "personal_extension/sdrms/routers"
	_ "personal_extension/sdrms/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	orm.RunSyncdb("default", false, true)
	beego.Run()
}
