package main

import (
	"personal_extension/Task/models"
	"personal_extension/Task/service"
	"personal_extension/lib/database"
	"personal_extension/lib/system"
)

func initSys() {
	config := system.Config{AppName: "", RunMode: system.DEV, ServerName: "", Filename: "app.conf", Dir: "./conf", AppConfigProvider: "ini"}
	system.InitSystem(&config)
	database.InitDatabase()
}

func main() {
	initSys()
	models.Register()
	service.ServerStart()
}
