package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ucpr/workspace2022/grpc_api_gateway/protos"
)

type Service struct {
}

func (s *Service) Get(ctx context.Context, in *pb.GetPingRequest) (*pb.GetPingResponse, error) {
	return &pb.GetPingResponse{
		Msg: "Pong",
	}, nil
}

func main() {
	s := grpc.NewServer()

	svc := &Service{}

	pb.RegisterServiceServer(s, svc)
	reflection.Register(s)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("not setting $PORT")
	}
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panicln("failed to listen port", err)
	}

	go func() {
		if err := s.Serve(l); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	s.GracefulStop()

	log.Println("finish server")
}
