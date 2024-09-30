package main

import (
	"context"
	pb "grps_log/pkg/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	GRPCADDR = "127.0.0.1:50051"
)

type Messager struct {
	Data []pb.Msg

	// Реализуем какой-то интерфейс из сгенерированного pb
	pb.UnimplementedMessagerServer
}

func (m *Messager) NewMessage(_ context.Context, msg *pb.Msg) (*pb.Empty, error) {
	m.Data = append(m.Data, *msg)
	return new(pb.Empty), nil
}

// Получать сообщения
func (m *Messager) Messages(_ *pb.Empty, stream pb.Messager_MessagesServer) error {
	for i := range m.Data {
		stream.Send(&m.Data[i])
	}
	return nil
}

func main() {
	srv := Messager{}

	lis, err := net.Listen("tcp4", GRPCADDR)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessagerServer(grpcServer, &srv)
	grpcServer.Serve(lis)
}
