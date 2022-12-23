package main

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn13/pb"
)

const (
	address     = "localhost:50051"
)

func main() {
	//实现方式1：
	//md := metadata.Pairs("key1", "value1", "key2", "key2")
	//ctx := metadata.NewOutgoingContext(context.Background(), md)
	// 实现方式2，设置两组K
	md := metadata.New(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 访问服务端address，创建连接conn
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	var header, trailer metadata.MD
	// 设置客户端访问超时时间1秒
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	// 客户端调用服务端 SayHello 请求，传入Name 为 "world", 返回值为服务端返回参数
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for key, value := range header {
		log.Printf("server header key:%v, value:%v\n", key, value)
	}
	for key, value := range trailer {
		log.Printf("server trailer key:%v, value:%v\n", key, value)
	}


	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r.GetMessage())
}