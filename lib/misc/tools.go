package misc

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

type value struct {
	val reflect.Value
	tag reflect.StructTag
	tpe reflect.Type
}

//返回两个结构体(必须有相同的内部变量名称和类型)中值不相同的字段, 对于b中未赋值的字段不作为比较的条件， 对于初始值为0 or ""的字段作为未赋值， 返回为true时表示有值不相同的字段，同时返回所有值不相同的字段名称, 添加exceptedFields将不作为比较的条件
func StructFac(a, b interface{}, exceptedFields ...string) (bool, []string) {
	at := GetType(a)
	av := GetValue(a)
	bt := GetType(b)
	bv := GetValue(b)

	aFieldNum := av.NumField()
	bFieldNum := bv.NumField()
	if aFieldNum != bFieldNum {
		return false, nil
	}
	aFieldMap := make(map[string]*value)
	bFieldMap := make(map[string]*value)
	for i := 0; i < aFieldNum; i++ {
		nv := new(value)
		cn := at.Field(i).Name
		if isExistAry(cn, exceptedFields) {
			continue
		}
		nv.tpe = at.Field(i).Type
		nv.val = av.Field(i)
		aFieldMap[at.Field(i).Name] = nv
	}
	for i := 0; i < bFieldNum; i++ {
		nv := new(value)
		nv.tpe = bt.Field(i).Type
		nv.val = bv.Field(i)
		bFieldMap[bt.Field(i).Name] = nv
	}
	unEqualFields := make([]string, 0)
	for k, v := range aFieldMap {
		if cv, ok := bFieldMap[k]; ok {
			if cv.tpe != v.tpe {
				return false, nil
			}
			bValid := bv.FieldByName(k)
			if !bValid.CanInterface() {
				continue
			}
			b := IsDefault(bValid.Interface())
			if cv.val.CanInterface() && v.val.CanInterface() {
				if !reflect.DeepEqual(cv.val.Interface(), v.val.Interface()) && !b {
					unEqualFields = append(unEqualFields, k)
				}
			}
		} else {
			return false, nil
		}
	}
	return true, unEqualFields
}

//判断变量声明后是否被赋值，其中默认值 将会当作未赋值做处理, 返回为ture时表示为默认值
func IsDefault(val interface{}) bool {
	b := true
	tpe := reflect.TypeOf(val)
	v := GetValue(val)
	if v.IsValid() {
		return isDefault(tpe, v, b)
	}
	return b
}

func isDefault(tpe reflect.Type, vv reflect.Value, status bool) bool {
	kd := tpe.Kind()
	switch kd {
	case reflect.Ptr:
		iv := reflect.Indirect(vv)
		if iv.IsValid() {
			b := isDefault(iv.Type(), iv, status)
			if b && status {
				status = true
			} else {
				status = false
			}
		} else {
			status = true
		}
	case reflect.Interface:
		var strDefault, integerDefault interface{}
		strDefault = ""
		integerDefault = 0
		if (reflect.DeepEqual(vv, strDefault)) || (reflect.DeepEqual(vv, integerDefault)) { //TODO
			status = true
		} else {
			status = false
		}
	case reflect.Map:
		kv := vv.MapKeys()
		for _, kf := range kv {
			b := isDefault(vv.MapIndex(kf).Type(), vv.MapIndex(kf), status)
			if b && status {
				status = true
			} else {
				status = false
			}
		}
	case reflect.Slice:
		for i := 0; i < vv.Len(); i++ {
			b := isDefault(vv.Index(i).Type(), vv.Index(i), status)
			if b && status {
				status = true
			} else {
				status = false
			}
		}
	case reflect.String:
		if vv.String() == "" {
			status = true
		} else {
			status = false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if vv.Uint() == 0 {
			status = true
		} else {
			status = false
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		if vv.Int() == 0 {
			status = true
		} else {
			status = false
		}
	case reflect.Float32, reflect.Float64:
		if vv.Float() == 0 {
			status = true
		} else {
			status = false
		}
	case reflect.Struct:
		if !vv.IsValid() {
			status = false //TODO
		}
		for i := 0; i < vv.NumField(); i++ {
			b := isDefault(vv.Field(i).Type(), vv.Field(i), status)
			if b && status {
				status = true
			} else {
				status = false
			}
		}
	}
	return status
}

func GetType(i interface{}) reflect.Type {
	t := reflect.TypeOf(i)
	for {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		} else {
			break
		}
	}
	return t
}

func GetValue(i interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(i))
}

func beginWithUpper(s string) bool {
	bt := []byte(s)
	if bt[0] >= 'A' && bt[0] <= 'Z' {
		return true
	}
	return false
}

func isExistAry(e interface{}, ary []interface{}) bool {
	if len(ary) == 0 {
		return false
	}
	ct := GetType(e)
	for _, v := range ary {
		t := GetType(v)
		if v == e && ct.Kind() == t.Kind() {
			return true
		}
		continue
	}
	return false
}

func Ips() []string {
	ips := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("addrs err: ", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

func GetHostByRoot(hostRoot string) string {
	ip := ""
	ips := Ips()
	if len(ips) > 0 {
		ip = IpFilter(hostRoot, ips)
	}
	return ip
}

func IpFilter(iproot string, ips []string) string {
	rs := strings.Split(iproot, ".")
	for _, ip := range ips {
		ipParts := strings.Split(ip, ".")
		if len(ipParts) < len(rs) {
			continue
		}

		find := true
		for i, part := range rs {
			if ipParts[i] != part {
				find = false
				break
			}
		}
		if find {
			return ip
		}
	}
	return ""
}
