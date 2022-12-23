package main

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	pb "learn16/pb"
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

func JWTInterceptor (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md,ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "miss metadata")
	}

	authType, ok := md["auth-type"]
	if !ok {
		return nil, status.Errorf(codes.Internal, "miss auth type")
	}

	if authType[0] == "no-jwt" {
		username, ok := md["username"]
		if !ok {
			return nil, status.Errorf(codes.Internal, "miss username")
		}
		password, ok := md["password"]
		if !ok {
			return nil, status.Errorf(codes.Internal, "miss password")
		}

		if username[0] != "zhangsan" ||  password[0] != "jay" {
			return nil, status.Errorf(codes.Unauthenticated, "auth fail")
		}

		token, err := CreateToken(username[0], "selfsecret")
		if err != nil {
			return nil, status.Errorf(codes.Internal, "token create err:" + err.Error())
		}

		header := metadata.New(map[string]string{
			"token": token})

		grpc.SetHeader(ctx, header)
		log.Printf("no-jwt auth \n")
	} else {
		auth, ok := md["authorization"]
		if !ok {
			return nil, status.Errorf(codes.Internal, "miss authorization")
		}

		uid, err := ParseToken(strings.TrimPrefix(auth[0], "Bearer "), "selfsecret")
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token auth fail")
		}

		if uid != "zhangsan" {
			return nil, status.Errorf(codes.Unauthenticated, "username auth fail")
		}
		log.Printf("jwt auth \n")
	}

	return handler(ctx, req)
}

func CreateToken(uid, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  uid,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(token string, secret string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if uid, ok := claim.Claims.(jwt.MapClaims)["uid"].(string); ok {
		return uid, nil
	}
	return "", fmt.Errorf("fail parse")
}

// main方法 函数开始执行的地方
func main() {
	// 调用标准库，监听50051端口的tcp连接
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//创建grpc服务
	s := grpc.NewServer(grpc.UnaryInterceptor(JWTInterceptor))
	//将server对象，也就是实现SayHello方法的对象，与grpc服务绑定
	pb.RegisterGreeterServer(s, &server{})
	// grpc服务开始接收访问50051端口的tcp连接数据
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}