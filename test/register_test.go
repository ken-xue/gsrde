package test

import (
	"google.golang.org/grpc"
	"gsrde/builder"
	pb "gsrde/test/hello"
	"log"
	"net"
	"testing"
)

func TestRegisterServiceA(t *testing.T) {
	var addr = "localhost:65535"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHiServiceServer(s, &pb.Server{})
	log.Printf("server listening at %v", lis.Addr())
	//注册到etcd
	cfg := builder.Config{
		Endpoints: []string{"192.168.100.24:2379"},
	}
	reg := builder.NewRegister(cfg)
	reg.Register(pb.HiService_ServiceDesc.ServiceName, addr)
	defer reg.UnRegister(pb.HiService_ServiceDesc.ServiceName)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func TestRegisterServiceB(t *testing.T) {
	var addr = "localhost:65533"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHiServiceServer(s, &pb.Server{})
	log.Printf("server listening at %v", lis.Addr())
	//注册到etcd
	cfg := builder.Config{
		Endpoints: []string{"192.168.100.24:2379"},
	}
	reg := builder.NewRegister(cfg)
	reg.Register(pb.HiService_ServiceDesc.ServiceName, addr)
	defer reg.UnRegister(pb.HiService_ServiceDesc.ServiceName)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
