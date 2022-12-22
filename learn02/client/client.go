package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
 	"google.golang.org/grpc/resolver"
	"learn02/pb"
)

//全局注册Scheme为myservice 的Resolver Build
func init() {
	resolver.Register(&myServiceBuilder{})
}
type myServiceBuilder struct{
}
func (*myServiceBuilder) Scheme() string {
	return "myservice"
}
type myServiceResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}
// 创建Resolver实例
func (*myServiceBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &myServiceResolver{
		target: target,
		cc:     cc,
	}
	r.start()
	return r, nil
}
// 根据target不同，解析出不同的端口
func (r *myServiceResolver) start() {
	if r.target.Endpoint == "abc" {
		r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: ":50051"}}})
	} else if r.target.Endpoint == "efg" {
		r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: ":50052"}}})
	}
}
//再次解析使用的解析方式不变
func (r *myServiceResolver) ResolveNow(o resolver.ResolveNowOptions) {
	r.start()
}
func (*myServiceResolver) Close() {}





const (
	address1     = "myservice:///abc"
	address2     = "myservice:///efg"
)

func main() {
	// 访问服务端address，创建连接conn，地址格式"myservice:///abc"
	conn, err := grpc.Dial(address1, grpc.WithInsecure(), grpc.WithBlock())
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

	//"myservice:///efg"
	conn2, err2 := grpc.Dial(address2, grpc.WithInsecure(), grpc.WithBlock())
	if err2 != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn2.Close()
	c2 := pb.NewGreeterClient(conn2)

	// 设置客户端访问超时时间1秒
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()
	// 客户端调用服务端 SayHello 请求，传入Name 为 "world", 返回值为服务端返回参数
	r2, err2 := c2.SayHello(ctx2, &pb.HelloRequest{Name: "world"})
	if err2 != nil {
		log.Fatalf("could not greet: %v", err2)
	}
	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r2.GetMessage())
}