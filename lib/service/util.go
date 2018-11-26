package service

import (
	"fmt"
)

func ServerName(prefix string, name string) string {
	servName := fmt.Sprintf("/%s/%s", prefix, name)
	return servName
}
