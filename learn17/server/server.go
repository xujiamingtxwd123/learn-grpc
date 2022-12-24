package main

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	pb "learn17/pb"
	"log"
	"net"
	"time"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedPersonServer
}
func (*server)SayPerson(stream pb.Person_SayPersonServer) error {
	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("err:" + err.Error())
				break
			}
			stream.Send(&pb.PersonReply{Message: req.Name + "-r"})
		}
		log.Println("go exit")
	}()

	time.Sleep(time.Hour)
	return nil
}

// main方法 函数开始执行的地方
func main() {
	cert, err := tls.LoadX509KeyPair("certs/server.crt","certs/server.key")
	if err != nil {
		panic(err)
	}
	certPool := x509.NewCertPool()
	ca,err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		panic(err)
	}
	certPool.AppendCertsFromPEM(ca)

	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certPool,
	})

	// 调用标准库，监听50051端口的tcp连接
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//创建grpc服务
	s := grpc.NewServer(grpc.Creds(cred))
	//将server对象，也就是实现SayHello方法的对象，与grpc服务绑定
	pb.RegisterPersonServer(s, &server{})
	// grpc服务开始接收访问50051端口的tcp连接数据
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}