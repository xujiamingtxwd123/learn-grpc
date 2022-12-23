package main

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn01/pb"
)

const (
	address     = "localhost:50051"
)

type ba struct {
	username string
	password string
}
func (b *ba) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	res := base64.StdEncoding.EncodeToString([]byte(b.username+":"+b.password))
	return map[string]string{"authorization": "Basic " + res }, nil
}

func (*ba) RequireTransportSecurity() bool {
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
	b := &ba {
		username: "zhangsan",
		password: "jay",
	}
	// 客户端调用服务端 SayHello 请求，传入Name 为 "world", 返回值为服务端返回参数
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"}, grpc.PerRPCCredentials(b))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r.GetMessage())
}