package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"gsrde/builder"
	"gsrde/discover"
	pb "gsrde/test/hello"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := builder.Config{
		Endpoints: []string{"192.168.100.24:2379"},
	}
	r := discover.NewResolver(cfg)
	resolver.Register(r)

	conn, err := grpc.Dial(r.Scheme()+"://author/"+pb.HiService_ServiceDesc.ServiceName,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHiServiceClient(conn)
	ticker := time.NewTicker(5 * time.Second)
	i := 1
	for range ticker.C {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		resp, err := c.SayHi(ctx, &pb.HiRequest{Name: "hi ,I am " + strconv.Itoa(i) + " message"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("message: %s", resp.GetMessage())
		i++
	}
}
