package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"learn03/pb"
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
	//模拟myservice解析出两个Address
	r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: ":50051"}, {Addr: ":50052"}}})
}
//再次解析使用的解析方式不变
func (r *myServiceResolver) ResolveNow(o resolver.ResolveNowOptions) {
	r.start()
}
func (*myServiceResolver) Close() {}

const (
	address1     = "myservice:///abc"
)



func init() {
	balancer.Register(newMyPickBuilder())
}

func newMyPickBuilder() balancer.Builder {
	return &myPickBuilder{}
}

type myPickBuilder struct{}

func (*myPickBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
	return &myPickBalancer{
		state:     0,
		cc:        cc,
		subConns:  make(map[resolver.Address]balancer.SubConn),
		subConns1: make(map[balancer.SubConn]resolver.Address),
	}
}

func (*myPickBuilder) Name() string {
	return "mypickBalance"
}



type myPickBalancer struct {
	state   connectivity.State
	cc    balancer.ClientConn
	subConns map[resolver.Address]balancer.SubConn
	subConns1 map[balancer.SubConn]resolver.Address
}

func (b *myPickBalancer) ResolverError(err error) {
	//TODO 需要剔除无效连接
}

func (b *myPickBalancer) UpdateClientConnState(s balancer.ClientConnState) error {
	addrsSet := make(map[resolver.Address]struct{})
	for _, a := range s.ResolverState.Addresses {
		addrsSet[a] = struct{}{}
		if _, ok := b.subConns[a]; !ok {
			sc, err := b.cc.NewSubConn([]resolver.Address{a}, balancer.NewSubConnOptions{})
			if err != nil {
				continue
			}
			b.subConns[a] = sc
			sc.Connect()
		}
	}
	return nil
}

func (b *myPickBalancer) UpdateSubConnState(sc balancer.SubConn, s balancer.SubConnState) {
	// TODO 需要剔除无效连接，增加有效连接
	if s.ConnectivityState == connectivity.Ready {
		b.subConns[b.subConns1[sc]] = sc
	}
	var scs []balancer.SubConn
	for _, sc := range b.subConns {
		scs = append(scs, sc)
	}

	if len(b.subConns) == 2{
		b.cc.UpdateState(balancer.State{ConnectivityState: connectivity.Ready, Picker: &myPicker{scs}})
	}
}

func (b *myPickBalancer) Close() {
}

type myPicker struct {
	subConns []balancer.SubConn
}

func (p *myPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	//获取当前时间
	second := time.Now().Second()
	fmt.Printf("Current Time Second:%d\n", second)
	if second % 2 == 0 {
		return balancer.PickResult{SubConn: p.subConns[0]}, nil
	}
	return balancer.PickResult{SubConn: p.subConns[1]}, nil
}




func main() {
	// 访问服务端address，创建连接conn，地址格式"myservice:///abc"
	conn, err := grpc.Dial(address1, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"mypickBalance":{}}]}`),)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
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

	time.Sleep(time.Second)
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()
	// 客户端调用服务端 SayHello 请求，传入Name 为 "world", 返回值为服务端返回参数
	r2, err2 := c.SayHello(ctx2, &pb.HelloRequest{Name: "world"})
	if err2 != nil {
		log.Fatalf("could not greet: %v", err2)
	}
	// 根据服务端处理逻辑，返回值也为"world"
	log.Printf("Greeting: %s", r2.GetMessage())
}