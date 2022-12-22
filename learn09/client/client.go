package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn09/pb"
)

const (
	address     = "localhost:50051"
)

func main() {
	// 访问服务端address，创建连接conn
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPersonClient(conn)

	// 设置客户端访问超时时间1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.SayPerson(ctx, &pb.PersonRequest{Name: "prefix"})
	if err != nil {
		log.Fatalf("could not stream: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		log.Println(res.Message)
	}
}