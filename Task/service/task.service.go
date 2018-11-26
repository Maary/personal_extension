package service

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"personal_extension/Task/models"
	"personal_extension/Task/task_service"
	"personal_extension/lib/misc"
	"personal_extension/lib/service"
	"personal_extension/lib/system"
	"time"
)

type TaskServer struct {
}

func (t TaskServer) QueryTasks(ctx context.Context, params *task_rpc_config.Params) (*task_rpc_config.Tasks, error) {
	spar := new(models.Task_QueryParam)
	err := json.Unmarshal([]byte(params.ParamsStr), spar)
	if err != nil {
		return nil, err
	}
	ts, err := models.QueryTasks(spar)
	if err != nil {
		return nil, err
	}
	jsonB, err := json.Marshal(ts)
	if err != nil {
		return nil, err
	}
	return &task_rpc_config.Tasks{Content: string(jsonB)}, nil
}

func (t TaskServer) InsertTasks(ctx context.Context, tasks *task_rpc_config.Tasks) (*task_rpc_config.SingleStatus, error) {
	tasksProto := make([]*models.Task, 0)
	err := json.Unmarshal([]byte(tasks.Content), &tasksProto)
	if err != nil {
		return nil, err
	}
	_, err = models.InsertTasks(tasksProto) //TODO: transport count to caller
	if err != nil {
		return nil, err
	}
	return &task_rpc_config.SingleStatus{""}, nil

}

func ServerStart() {
	host := misc.GetHostByRoot(system.AppConfig.String("host_root"))
	port := system.AppConfig.String("port")
	etcdHosts := system.AppConfig.Strings("etcd::hosts")
	prefix := system.AppConfig.String("etcd::prefix")
	serviceName := system.AppConfig.String("etcd::task_server")
	fmt.Println(host, port, etcdHosts, prefix, serviceName)
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("Failed to listen tcp: %v", err)
	}
	if len(etcdHosts) > 0 && prefix != "" {
		err = service.InitEtcd(etcdHosts, prefix)
		if err != nil {
			log.Fatalf("Init etcd for service ERROR: %v", err)
		}

		startErrOrNo := service.Start(serviceName, host, port, time.Second*10, 15)
		if startErrOrNo != nil {
			log.Fatalf("faield to start service: %v", err)
		} else {
			fmt.Println("service register is ok!")
		}
	}
	s := grpc.NewServer()
	task_rpc_config.RegisterTaskServer(s, &TaskServer{})
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Start server ERROR: %v", err)
	}
}
