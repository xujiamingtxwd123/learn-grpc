package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn10/pb"
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

	// 设置客户端访问超时时间100秒
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	stream, err := c.SayPerson(ctx)
	if err != nil {
		log.Fatalf("could not stream: %v", err)
	}

	go func() {
		stream.Send(&pb.PersonRequest{Name: "go-1-1"})
		time.Sleep(100 * time.Millisecond)
		stream.Send(&pb.PersonRequest{Name: "go-1-2"})
		time.Sleep(100 * time.Millisecond)
		stream.Send(&pb.PersonRequest{Name: "go-1-3"})
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil{
				log.Println("err:" + err.Error())
				break
			}
			log.Println(res.Message)
		}
	}()
	log.Println("over")
	time.Sleep(time.Hour)
}