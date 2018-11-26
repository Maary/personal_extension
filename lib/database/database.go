package database

import (
	"crawler.center/lib/system"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() {
	//读取配置文件，设置数据库参数
	//数据库类别
	dbType := system.AppConfig.String("db_type")
	//连接名称
	dbAlias := system.AppConfig.String(dbType + "::db_alias")
	//数据库名称
	dbName := system.AppConfig.String(dbType + "::db_name")
	//数据库连接用户名
	dbUser := system.AppConfig.String(dbType + "::db_user")
	//数据库连接用户名
	dbPwd := system.AppConfig.String(dbType + "::db_pwd")
	//数据库IP（域名）
	dbHost := system.AppConfig.String(dbType + "::db_host")
	//数据库端口
	dbPort := system.AppConfig.String(dbType + "::db_port")
	switch dbType {
	case "sqlite3":
		orm.RegisterDataBase(dbAlias, dbType, dbName)
	case "mysql":
		dbCharset := system.AppConfig.String(dbType + "::db_charset")
		dataSouce := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Asia%%2FShanghai", dbUser, dbPwd, dbHost, dbPort, dbName, dbCharset)
		orm.RegisterDataBase(dbAlias, dbType, dataSouce, 30)
	}
	//如果是开发模式，则显示命令信息
	isDev := (system.AppConfig.String("runmode") == "dev")
	if isDev {
		orm.Debug = isDev
	}
}
