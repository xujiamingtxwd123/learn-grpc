package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "learn13/pb"
	"log"
	"net"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// 该函数定义必须与helloworld.pb.go 定义的SayHello一致
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	//打印客户端传入HelloRequest请求的Name参数
	log.Printf("Received: %v", in.GetName())

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Fatal("FromIncomingContext error")
	}
	for key, value := range md {
		log.Printf("key:%v, value:%v\n", key, value)
	}

	//设置 服务端HEADERS帧 headerfields
	header := metadata.New(map[string]string{
		"ResKey1": "ResValue1",
		"ResKey2": "ResVaule2"})
	grpc.SetHeader(ctx, header)
	//可以调用SendHeader 手动发送 HEADERS
	//grpc.SendHeader()

	trailer := metadata.New(map[string]string{
		"TrResKey1": "ResValue1",
		"TrResKey2": "ResVaule2"})
	grpc.SetTrailer(ctx, trailer)

	//将name参数作为返回值，返回给客户端
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// main方法 函数开始执行的地方
func main() {
	// 调用标准库，监听50051端口的tcp连接
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//创建grpc服务
	s := grpc.NewServer()
	//将server对象，也就是实现SayHello方法的对象，与grpc服务绑定
	pb.RegisterGreeterServer(s, &server{})
	// grpc服务开始接收访问50051端口的tcp连接数据
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}