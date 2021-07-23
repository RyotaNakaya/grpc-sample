// Package main implements a server for Greeter service.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "github.com/RyotaNakaya/grpc-sample/helloworld_stream/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) SayHelloStream(stream pb.Greeter_SayHelloStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{
				Message: "done",
			})
		}
		if err != nil {
			return err
		}
		fmt.Println(req.GetName())
	}
}

func (s *server) SayHelloListStream(stream pb.Greeter_SayHelloListStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ListHelloReply{
				Message: "done",
			})
		}
		if err != nil {
			return err
		}

		var filename string = "./helloworld_stream/server/list.txt"
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// slice が流れてくるので、ループで一行ずつ書き出す
		for _, v := range req.NameList {
			_, err := file.WriteString(fmt.Sprintf("%s\n", v))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
