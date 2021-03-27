package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/jefersonvinicius/microservices-course-communication/grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	AddUserVerbose(client)
	// AddUser(client)

}

func AddUser(client pb.UserServiceClient) { // Using request unidirectional

	req := &pb.User{
		Id:    "0",
		Name:  "Jeferson Vinícius",
		Email: "j@gmail.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) { // Using request with stream
	req := &pb.User{
		Id:    "0",
		Name:  "Jeferson Vinícius",
		Email: "j@gmail.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}

		fmt.Println("Status:", stream.Status)
	}
}
