package main

import (
	"context"
	"google.golang.org/grpc/encoding"
	"io"
	"log"
	"time"

	"github.com/pierrec/lz4/v4"
	"google.golang.org/grpc"
	"learn14/pb"
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
	c := pb.NewGreeterClient(conn)

	// 设置客户端访问超时时间1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 客户端调用服务端 SayHello 请求，传入Name 为 "world", 返回值为服务端返回参数
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"}, grpc.UseCompressor("lz4"))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r.GetMessage())
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