package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn08/pb"
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

	stream, err := c.SayPerson(ctx)
	if err != nil {
		log.Fatalf("could not stream: %v", err)
	}

	stream.Send(&pb.PersonRequest{Name: "zhangsan"})
	stream.Send(&pb.PersonRequest{Name: "lisi"})
	stream.Send(&pb.PersonRequest{Name: "wangwu"})

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("stream recv: %v", err)
	}

	log.Printf("count: %d", res.GetCount())
}