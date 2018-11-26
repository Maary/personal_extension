package service

import (
	"errors"
	"fmt"

	etcd3 "github.com/coreos/etcd/clientv3"
)

var client *etcd3.Client
var prefix string

func InitEtcd(hosts []string, prefixN string) error {
	var err error
	if prefixN == "" {
		return errors.New("Etcd prefix can't be null")
	}
	prefix = prefixN
	client, err = etcd3.New(etcd3.Config{Endpoints: hosts})
	if err != nil {
		return fmt.Errorf("grpclb: create etcd3 client failed: %v", err)
	}
	return nil
}
