// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/RyotaNakaya/grpc-sample/helloworld_stream/proto"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	res, err := request(c)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("result: %v", res)
}

func request(client pb.GreeterClient) (*pb.HelloReply, error) {
	stream, err := client.SayHelloStream(context.Background())
	if err != nil {
		return nil, err
	}

	const loop = 10
	for i := 0; i <= loop; i++ {
		r := &pb.HelloRequest{
			Name: fmt.Sprintf("Taro%d", i),
		}
		if err := stream.Send(r); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		time.Sleep(time.Second * 1)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return res, nil
}
