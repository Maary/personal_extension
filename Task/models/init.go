package models

import (
	"personal_extension/lib/system"
	"personal_extension/lib/task"
)

func Register() {
	task.Register(new(Task))
}

func TableName(name string) string {
	prefix := system.AppConfig.String("db_dt_prefix")
	return prefix+name
}