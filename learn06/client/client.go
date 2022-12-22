package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn06/pb"
)

const (
	address     = "localhost:50051"
)

var (
	retryPolicy = `{
		"RetryThrottling": {
		  "MaxTokens": 4,
		  "TokenRatio": 0.1
		},
		"MethodConfig": [{
		"Name": [{"Service": "Greeter"}],
		  "RetryPolicy": {
			  "MaxAttempts": 6,
			  "InitialBackoff": "2s",
			  "MaxBackoff": "10s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)
// "Service": "" 表示全局应用
func main() {
	// 访问服务端address，创建连接conn
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultServiceConfig(retryPolicy))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 设置客户端访问超时时间1秒
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()
	// 客户端调用服务端 SayHello 请求，传入Name 为 "world", 返回值为服务端返回参数
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r.GetMessage())
}