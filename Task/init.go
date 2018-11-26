package Task

import (
	"github.com/astaxie/beego/orm"
	"personal_extension/lib/system"
)

func Init() {
	orm.RegisterModel(new(Result))
}

func TableName(name string) string {
	prefix := system.AppConfig.String("")
	return prefix+name
}