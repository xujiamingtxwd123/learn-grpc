package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn11/pb"
)

const (
	address     = "localhost:50051"
)

func TimeInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	now := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("method:%s, exec:%d", method, time.Now().Sub(now).Milliseconds())
	return err
}

func main() {
	// 访问服务端address，创建连接conn
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithChainUnaryInterceptor(TimeInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 设置客户端访问超时时间1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 客户端调用服务端 SayHello 请求，传入Name 为 "world", 返回值为服务端返回参数
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r.GetMessage())
}