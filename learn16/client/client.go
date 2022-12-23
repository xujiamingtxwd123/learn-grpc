package main

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn16/pb"
)

const (
	address     = "localhost:50051"
)

type tokenAuth struct {
	username string
	password string
	token 	string
}
func (b *tokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if b.token != "" {
		return map[string]string{
			"authorization": "Bearer " + b.token,
			"auth-type": "jwt",
		}, nil
	}

	return map[string]string{
		"username": b.username,
		"password": b.password,
		"auth-type": "no-jwt",
	}, nil
}

func (*tokenAuth) RequireTransportSecurity() bool {
	return false
}



func main() {
	// 访问服务端address，创建连接conn
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 设置客户端访问超时时间1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	b := &tokenAuth {
		username: "zhangsan",
		password: "jay",
	}
	var header metadata.MD
	// 客户端调用服务端 SayHello 请求，传入Name 为 "world", 返回值为服务端返回参数
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"},
	grpc.PerRPCCredentials(b), grpc.Header(&header))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r.GetMessage())


	token := header.Get("token")
	if token != nil {
		b.token = token[0]
	}

	//再次发起请求
	r1, err1 := c.SayHello(ctx, &pb.HelloRequest{Name: "world"},
		grpc.PerRPCCredentials(b), grpc.Header(&header))
	if err1 != nil {
		log.Fatalf("could not greet: %v", err1)
	}
	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r1.GetMessage())



}