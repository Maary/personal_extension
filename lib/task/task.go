package task

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"personal_extension/lib/misc"
	"reflect"
)

var tcc taskCache
func init() {
	tcc = make(taskCache)
}

type taskCache map[string]*taskProto

type structInfo struct {
	fieldTpe reflect.Type
	fieldTag reflect.StructTag
}

type taskProto struct {
	fieldInfo map[string]structInfo
	tpeName string
	content interface{}
}

func Register(tasks ...interface{}) {
	for _, task := range tasks {
		isValid(task)
		currentTp := misc.GetType(task)
		currentVal := misc.GetValue(task)
		fullName := fmt.Sprintf("%s/%s", currentTp.PkgPath(), currentTp.Name())
		ntp := new(taskProto)
		fieldsNum := currentTp.NumField()
		for i := 0; i < fieldsNum; i++ {
			n := new(structInfo)
			n.fieldTpe = currentTp.Field(i).Type
			n.fieldTag = currentTp.Field(i).Tag
		}
		name := getTableName(currentVal)
		ntp.tpeName = name
		ntp.content = task
		tcc[fullName] = ntp
		orm.RegisterModel(task)
	}
}

func isValid(task interface{}) {
	val := misc.GetValue(task)
	if !val.FieldByName("Type").IsValid() {
		panic("Task must have the unique [Type] field for business")
	}
	if !val.FieldByName("Id").IsValid() {
		panic("Task must have the [Id] field for storing database")
	}
	if !val.FieldByName("UUID").IsValid() && !val.FieldByName("Uuid").IsValid() {
		panic("Task must have the [Uuid] or [UUID] field for tracing the task path")
	}
	//if !val.MethodByName("TableName").IsValid() {
	//	panic("Task must have the [TableName] function for register the database model")
	//}
}

func getTableName(val reflect.Value) string {
	if fun := val.MethodByName("TableName"); fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		// has return and the first val is string
		if len(vals) > 0 && vals[0].Kind() == reflect.String {
			return vals[0].String()
		}
	}
	return misc.SnakeString(reflect.Indirect(val).Type().Name())
}

