package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn04/pb"
)


var kacp = keepalive.ClientParameters{
	Time:                15 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

const (
	address     = "localhost:50051"
)

func main() {
	// 访问服务端address，创建连接conn
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithKeepaliveParams(kacp),
		grpc.WithStatsHandler(&StatsHandler{}))
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


type StatsHandler struct {

}

//TagConn可以将一些信息附加到给定的上下文。
func (h *StatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	fmt.Printf("TagConn:%v\n", info)
	return ctx
}

// 会在连接开始和结束时被调用，分别会输入不同的状态.

func (h *StatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	fmt.Printf("HandleConn:%v\n", s)
}


// TagRPC可以将一些信息附加到给定的上下文

func (h *StatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	fmt.Printf("TagRPC:%v\n", info)
	return ctx
}

// 处理RPC统计信息
func (h *StatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	fmt.Printf("HandleRPC:%v\n", s)
}