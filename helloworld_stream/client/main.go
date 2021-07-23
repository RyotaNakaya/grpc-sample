// Package main implements a client for Greeter service.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
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
	// res, err := request(c)
	res, err := requestWithFile(c)
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

func requestWithFile(client pb.GreeterClient) (*pb.ListHelloReply, error) {
	fmt.Println("call requestWithFile")
	stream, err := client.SayHelloListStream(context.Background())
	if err != nil {
		return nil, err
	}

	const filename string = "./helloworld_stream/client/list.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var dump []string
	sliceLimit := 1000

	scanner := bufio.NewScanner(file)
	// ファイルから一行ずつ読み込む
	for scanner.Scan() {
		dump = append(dump, scanner.Text())
		// リミットに達したらストリームに流し込んで、slice を初期化
		if len(dump) == sliceLimit {
			r := &pb.ListHelloRequest{
				NameList: dump,
			}
			fmt.Println("send stream")
			if err := stream.Send(r); err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			dump = nil
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return res, nil
}
