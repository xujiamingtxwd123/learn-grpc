package main

import (
	"context"
	"github.com/pierrec/lz4/v4"
	"google.golang.org/grpc/encoding"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "learn14/pb"
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

func init() {
	c := &lz4compressor{}
	encoding.RegisterCompressor(c)
}

type lz4compressor struct {

}

func (c *lz4compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	writer := lz4.NewWriter(w)
	writer.Reset(w)
	return writer, nil
}

func (c *lz4compressor) Decompress(r io.Reader) (io.Reader, error) {
	return lz4.NewReader(r), nil
}

func (c *lz4compressor) Name() string {
	return "lz4"
}