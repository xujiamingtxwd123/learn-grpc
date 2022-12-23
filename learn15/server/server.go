package main

import (
	"context"
	"encoding/base64"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	pb "learn01/pb"
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
	//将name参数作为返回值，返回给客户端
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func BAInterceptor (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md,ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "miss metadata")
	}
	auth, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Internal, "invalid metadata")
	}
	r := base64.StdEncoding.EncodeToString([]byte("zhangsan:jay"))
	r1 := strings.TrimPrefix(auth[0], "Basic ")
	if r != r1 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth")
	}
	return handler(ctx, req)
}
// main方法 函数开始执行的地方
func main() {
	// 调用标准库，监听50051端口的tcp连接
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//创建grpc服务
	s := grpc.NewServer(grpc.UnaryInterceptor(BAInterceptor))
	//将server对象，也就是实现SayHello方法的对象，与grpc服务绑定
	pb.RegisterGreeterServer(s, &server{})
	// grpc服务开始接收访问50051端口的tcp连接数据
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}