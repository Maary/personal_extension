package service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"encoding/json"
	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

var stopSignal = make(chan bool, 1)

func Start(name string, host string, port string, interval time.Duration, ttl int) error {
	err := Register(name, host, port, interval, ttl)
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		UnRegister(name, host, port)
		os.Exit(1)
	}()

	return nil
}

// Register
func Register(name string, host string, port string, interval time.Duration, ttl int) error {
	serviceValue := fmt.Sprintf("%s:%s", host, port)
	serviceKey := ServerName(prefix, name)
	serviceKey = fmt.Sprintf("%s/%s", serviceKey, serviceValue)

	update := naming.Update{}
	update.Addr = serviceValue
	update.Op = naming.Add
	updateB, _ := json.Marshal(update)
	updateS := string(updateB)
	go func() {
		// invoke self-register with ticker
		ticker := time.NewTicker(interval)
		for {
			// minimum lease TTL is ttl-second
			resp, _ := client.Grant(context.TODO(), int64(ttl))
			// should get first, if not exist, set it
			_, err := client.Get(context.Background(), serviceKey)

			if err != nil {
				if err == rpctypes.ErrKeyNotFound {
					if _, err := client.Put(context.TODO(), serviceKey, updateS, etcd3.WithLease(resp.ID)); err != nil {
						log.Printf("grpclb: set service '%s' with ttl to etcd3 failed: %s", name, err.Error())
					}
				} else {
					log.Printf("grpclb: service '%s' connect to etcd3 failed: %s", name, err.Error())
				}
			} else {
				// refresh set to true for not notifying the watcher
				if _, err := client.Put(context.Background(), serviceKey, updateS, etcd3.WithLease(resp.ID)); err != nil {
					log.Printf("grpclb: refresh service '%s' with ttl to etcd3 failed: %s", name, err.Error())
				}
			}
			select {
			case <-stopSignal:
				return
			case <-ticker.C:
			}
		}
	}()

	return nil
}

// UnRegister delete registered service from etcd
func UnRegister(name string, host string, port string) error {
	serviceValue := fmt.Sprintf("%s:%s", host, port)
	serviceKey := ServerName(prefix, name)
	serviceKey = fmt.Sprintf("%s%s", serviceKey, serviceValue)

	stopSignal <- true
	stopSignal = make(chan bool, 1) // just a hack to avoid multi UnRegister deadlock
	var err error
	if _, err := client.Delete(context.Background(), serviceKey); err != nil {
		log.Printf("grpclb: deregister '%s' failed: %s", serviceKey, err.Error())
	} else {
		log.Printf("grpclb: deregister '%s' ok.", serviceKey)
	}
	return err
}
