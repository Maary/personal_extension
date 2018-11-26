package system

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/utils"
	"os"
	"path/filepath"
	"reflect"
)

type Config struct {
	AppName           string
	RunMode           string
	ServerName        string
	Filename          string // 配置文件名称
	Dir               string // 配置文件目录
	AppConfigProvider string // 配置文件类型

}

var (
	// 默认配置
	BConfig *Config

	// 应用程序配置
	AppConfig *defaultAppConfig

	IsInit bool = false
)

func InitSystem(c *Config) {
	if IsInit && c == nil {
		return
	}
	IsInit = true
	var err error
	if c == nil {
		BConfig = newBConfig()
	} else {
		BConfig = c
	}
	appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfigPath := filepath.Join(workPath, BConfig.Dir, BConfig.Filename)
	if !utils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(appPath, BConfig.Dir, BConfig.Filename)
		if !utils.FileExists(appConfigPath) {
			AppConfig = &defaultAppConfig{innerConfig: config.NewFakeConfig()}
			return
		}
	}
	if err = parseConfig(appConfigPath); err != nil {
		panic(err)
	}
}

// 默认配置
func newBConfig() *Config {
	return &Config{
		AppName:           "unknown",
		RunMode:           DEV,
		ServerName:        "hello word",
		Filename:          "app.conf",
		Dir:               "conf",
		AppConfigProvider: "ini",
	}
}

// 支持ini和json两种格式
func parseConfig(appConfigPath string) (err error) {
	AppConfig, err = newAppConfig(BConfig.AppConfigProvider, appConfigPath)
	if err != nil {
		return err
	}
	return assignConfig(AppConfig)
}

func newAppConfig(appConfigProvider, appConfigPath string) (*defaultAppConfig, error) {
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &defaultAppConfig{ac}, nil
}

func assignConfig(ac config.Configer) error {
	//todo 这里可以附加上rpc配置
	for _, i := range []interface{}{BConfig} {
		assignSingleConfig(i, ac)
	}
	// set the run mode first
	if runMode := ac.String("RunMode"); runMode != "" {
		BConfig.RunMode = runMode
	}
	return nil
}

func assignSingleConfig(p interface{}, ac config.Configer) {
	pt := reflect.TypeOf(p)
	if pt.Kind() != reflect.Ptr {
		return
	}
	pt = pt.Elem()
	if pt.Kind() != reflect.Struct {
		return
	}
	pv := reflect.ValueOf(p).Elem()

	for i := 0; i < pt.NumField(); i++ {
		pf := pv.Field(i)
		if !pf.CanSet() {
			continue
		}
		name := pt.Field(i).Name
		switch pf.Kind() {
		case reflect.String:
			pf.SetString(ac.DefaultString(name, pf.String()))
		case reflect.Int, reflect.Int64:
			pf.SetInt(ac.DefaultInt64(name, pf.Int()))
		case reflect.Bool:
			pf.SetBool(ac.DefaultBool(name, pf.Bool()))
		case reflect.Struct:
		default:
			//do nothing here
		}
	}
}
