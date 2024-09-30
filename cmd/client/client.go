package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "grps_log/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	CONADDR = "127.0.0.1:50051"
)

func main() {
	ctx := context.Background()

	conn, err := grpc.NewClient(CONADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMessagerClient(conn)

	start := time.Now()

	var i int
	for i := 0; i < 10; i++ {
		msg := pb.Msg{
			Id:   int64(i),
			Text: fmt.Sprintf("New message - %d", i),
		}
		client.NewMessage(ctx, &msg)
	}

	duration := time.Since(start)
	fmt.Printf("Отправлено %d сообщений за %v.\n", i, duration)

	start = time.Now()

	err = PrintAllMessage(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	duration = time.Since(start)
	fmt.Printf("Все сообщения приняты за %v.", duration)
}

func PrintAllMessage(ctx context.Context, client pb.MessagerClient) error {
	fmt.Println("Запрашиваем сообщения на севрере...")
	stream, err := client.Messages(ctx, &pb.Empty{})
	if err != nil {
		return err
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			i++
			fmt.Printf("Сообщение #%d: %v\n", i, msg.Text)
		}
	}
}
