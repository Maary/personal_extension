package service

import (
	"fmt"
	"time"

	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func Dial(servName string) (conn *grpc.ClientConn, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if prefix != "" {
		fullName := ServerName(prefix, servName)
		r := &etcdnaming.GRPCResolver{Client: client}
		b := grpc.RoundRobin(r)
		grpcConn, gErr := grpc.DialContext(ctx, fullName, grpc.WithInsecure(), grpc.WithBalancer(b))
		if gErr != nil {
			return nil, gErr
		}
		return grpcConn, nil
	}
	grpcConn, gErr := grpc.DialContext(ctx, servName, grpc.WithInsecure())
	if gErr != nil {
		fmt.Println("rpc error: ", gErr)
		return nil, gErr
	}
	return grpcConn, nil
}
