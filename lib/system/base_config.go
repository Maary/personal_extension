package system

import (
	"github.com/astaxie/beego/config"
)

type defaultAppConfig struct {
	innerConfig config.Configer
}

func (b *defaultAppConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(BConfig.RunMode+"::"+key, val); err != nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *defaultAppConfig) String(key string) string {
	if v := b.innerConfig.String(BConfig.RunMode + "::" + key); v != "" {
		return v
	}
	return b.innerConfig.String(key)
}

func (b *defaultAppConfig) Strings(key string) []string {
	if v := b.innerConfig.Strings(BConfig.RunMode + "::" + key); len(v) > 0 {
		return v
	}
	return b.innerConfig.Strings(key)
}

func (b *defaultAppConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(key)
}

func (b *defaultAppConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(key)
}

func (b *defaultAppConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(key)
}

func (b *defaultAppConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(BConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(key)
}

func (b *defaultAppConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *defaultAppConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *defaultAppConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *defaultAppConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *defaultAppConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *defaultAppConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *defaultAppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *defaultAppConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *defaultAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}
