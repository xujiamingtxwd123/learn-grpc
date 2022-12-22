package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	pb "learn09/pb"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedPersonServer
}
func (*server)SayPerson(req *pb.PersonRequest, stream pb.Person_SayPersonServer) error {
	stream.Send(&pb.PersonReply{
		Message: req.Name + "-1",
	})
	stream.Send(&pb.PersonReply{
		Message: req.Name + "-2",
	})
	stream.Send(&pb.PersonReply{
		Message: req.Name + "-3",
	})
	return nil
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
	pb.RegisterPersonServer(s, &server{})
	// grpc服务开始接收访问50051端口的tcp连接数据
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}