package main

import (
	"context"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/reflect/protoreflect"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	"log"
)

const (
	address     = "localhost:50051"
)

func main() {
	// 访问服务端address，创建连接conn
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	client, err := c.ServerReflectionInfo(context.Background())
	if err != nil {
		log.Fatalf("ServerReflectionInfo err:%v", err)
		return
	}
	// 1. 获取Service List
	req :=  grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices {},
	}
	client.Send(&req)
	r, err := client.Recv()
	if err != nil {
		log.Fatalf("recv: %v", err)
		return
	}

	log.Printf("ServiceList:%v", r.GetListServicesResponse().String())

	//2.获取文件描述符
	req =  grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol {
			FileContainingSymbol: "Greeter",
		},
	}
	client.Send(&req)
	r, err = client.Recv()
	if err != nil {
		log.Fatalf("recv: %v", err)
		return
	}

	fd := new(descriptorpb.FileDescriptorProto)
	if err := proto.Unmarshal(r.GetFileDescriptorResponse().GetFileDescriptorProto()[0], fd); err != nil {
		log.Fatalf("unmarshal file descriptor : %v", err)
		return
	}
	log.Printf("FileDescriptor:%v", fd.String())

	// 3. 方法调用
	helloRequestReflect := fd.MessageType[0].ProtoReflect()
	helloRequestDescriptor := helloRequestReflect.Descriptor()
	//也就是HelloRequest Name字段 proto文件该字段序号为1
	firstArgsFieldDescriptor := helloRequestDescriptor.Fields().ByNumber(1)
	helloRequestReflect.Set(firstArgsFieldDescriptor, protoreflect.ValueOfString("zhangsan"))

	args := helloRequestReflect.Interface()

	ret := fd.MessageType[1].ProtoReflect().Interface()

	err = conn.Invoke(context.Background(), "/Greeter/SayHello", args, ret)
	if err != nil {
		log.Fatalf("invoke: %v", err)
		return
	}
	log.Println(ret)
}