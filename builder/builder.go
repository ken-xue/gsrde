package builder

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"gsrde/register"
	"log"
	"time"
)

type Config struct {
	Endpoints []string
	Timeout   time.Duration
}

func NewRegister(cfg Config) register.DefaultRegister {
	register := register.DefaultRegister{}
	var err error
	register.Cli, err = NewEtcdClient(cfg)
	if err != nil {
		log.Fatalf("error : %v", err)
	}
	return register
}

func NewEtcdClient(cfg Config) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.Timeout,
	})
	if err != nil {
		log.Fatalf("error : %v", err)
	}
	return cli, err
}
