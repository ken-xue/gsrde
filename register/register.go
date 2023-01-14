package register

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type Register interface {
	Register(service, addr string) error
	UnRegister(service string) error
}

type DefaultRegister struct {
	Cli *clientv3.Client
}

// Register 注册一个rpc服务
func (r *DefaultRegister) Register(serviceName, addr string) (err error) {
	key := fmt.Sprintf("%s/%s", serviceName, addr)
	//创建租约
	leaseResp, err := r.Cli.Grant(context.Background(), 10)
	if err != nil {
		log.Fatalf("Create Lease error ：%s\n", err)
		return err
	}
	//注册到etcd
	log.Printf("Register to etcd : key : %s\n", key)
	_, err = r.Cli.Put(context.Background(), key, key, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Fatalf("Register to etcd error : %s\n", err)
	}
	//开启心跳检查
	ch, err := r.Cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		log.Fatalf("KeepAlive to etcd error : %s\n", err)
	}
	//清空ch的数据
	go func() {
		for {
			<-ch
		}
	}()
	return
}

// UnRegister 取消注册一个rpc服务
func (r *DefaultRegister) UnRegister(service string) error {
	resp, err := r.Cli.Delete(context.Background(), service)
	if err != nil {
		log.Fatalf("UnRegister error : %s, service : %v, resp : %v\n", err, service, resp)
	}
	return err
}
