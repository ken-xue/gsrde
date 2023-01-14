package hello

import (
	"context"
	"log"
)

type Server struct {
	UnimplementedHiServiceServer
}

func (s *Server) SayHi(ctx context.Context, in *HiRequest) (*HiReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &HiReply{Message: "Hi " + in.GetName()}, nil
}
