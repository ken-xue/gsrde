package discover

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"gsrde/builder"
	"log"
	"strings"
)

// EtcdResolver etcd解析器
type EtcdResolver struct {
	cli  *clientv3.Client
	conn resolver.ClientConn
}

// NewResolver 初始化一个etcd解析器
func NewResolver(cfg builder.Config) resolver.Builder {
	cli, err := builder.NewEtcdClient(cfg)
	if err != nil {
		log.Fatalf("error : %v", err)
	}
	return &EtcdResolver{
		cli: cli,
	}
}

func (r *EtcdResolver) Scheme() string {
	return "etcd"
}

func (r *EtcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

func (r *EtcdResolver) Close() {
}

// Build 构建解析器
func (r *EtcdResolver) Build(target resolver.Target, ClientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.conn = ClientConn
	// 监听key的变化
	go r.watch(fmt.Sprintf("%s/", target.Endpoint))
	return r, nil
}

// 监听etcd中某个key前缀的服务地址列表的变化
func (r *EtcdResolver) watch(keyPrefix string) {
	//初始化服务地址列表
	var addresses []resolver.Address
	resp, err := r.cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("Get service error ：", err)
	} else {
		for i := range resp.Kvs {
			addresses = append(addresses,
				resolver.Address{
					Addr: strings.TrimPrefix(string(resp.Kvs[i].Key), keyPrefix),
				},
			)
		}
	}
	status := resolver.State{
		Addresses: addresses,
	}
	r.conn.UpdateState(status)
	//监听服务地址列表的变化
	rch := r.cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, event := range n.Events {
			addr := strings.TrimPrefix(string(event.Kv.Key), keyPrefix)
			switch event.Type {
			case mvccpb.PUT:
				if !exists(addresses, addr) {
					addresses = append(addresses, resolver.Address{Addr: addr})
					status.Addresses = addresses
					r.conn.UpdateState(status)
				}
				log.Printf("service register ：%s", addr)
			case mvccpb.DELETE:
				if s, ok := remove(addresses, addr); ok {
					status.Addresses = s
					r.conn.UpdateState(status)
				}
				log.Printf("service destroy ：%s", addr)
			}
		}
	}
}

func exists(addresses []resolver.Address, addr string) bool {
	for i := range addresses {
		if addresses[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(addresses []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range addresses {
		if addresses[i].Addr == addr {
			addresses[i] = addresses[len(addresses)-1]
			return addresses[:len(addresses)-1], true
		}
	}
	return nil, false
}
