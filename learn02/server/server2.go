package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "learn02/pb"
)

const (
	port2 = ":50052"
)

type server2 struct {
	pb.UnimplementedGreeterServer
}

// 该函数定义必须与helloworld.pb.go 定义的SayHello一致
func (s *server2) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	//打印客户端传入HelloRequest请求的Name参数
	log.Printf("Received: %v", in.GetName())
	//将name参数作为返回值，返回给客户端
	return &pb.HelloReply{Message: "Service2: Hello " + in.GetName()}, nil
}

// main方法 函数开始执行的地方
func main() {
	// 调用标准库，监听50052端口的tcp连接
	lis, err := net.Listen("tcp", port2)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//创建grpc服务
	s := grpc.NewServer()
	//将server对象，也就是实现SayHello方法的对象，与grpc服务绑定
	pb.RegisterGreeterServer(s, &server2{})
	// grpc服务开始接收访问50051端口的tcp连接数据
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}