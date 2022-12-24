package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	"learn17/pb"
)

const (
	address     = "localhost:50051"
)


func main() {
	cert, err := tls.LoadX509KeyPair("certs/client.crt","certs/client.key")
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
		ServerName: "myserver.com",
		RootCAs: certPool,
	})

	// 访问服务端address，创建连接conn
	conn, err := grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(cred), grpc.WithTimeout(10 * time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPersonClient(conn)

	// 设置客户端访问超时时间100秒
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	stream, err := c.SayPerson(ctx)
	if err != nil {
		log.Fatalf("could not stream: %v", err)
	}

	go func() {
		stream.Send(&pb.PersonRequest{Name: "go-1-1"})
		time.Sleep(100 * time.Millisecond)
		stream.Send(&pb.PersonRequest{Name: "go-1-2"})
		time.Sleep(100 * time.Millisecond)
		stream.Send(&pb.PersonRequest{Name: "go-1-3"})
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil{
				log.Println("err:" + err.Error())
				break
			}
			log.Println(res.Message)
		}
	}()
	log.Println("over")
	time.Sleep(time.Hour)
}