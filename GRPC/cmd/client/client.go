package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
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

func AddUserVerbose(client pb.UserServiceClient) { // Getting data with stream
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

func AddUsers(client pb.UserServiceClient) { // Sending data with stream
	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "Jeferson",
			Email: "jeferson@gmail.com",
		},
		{
			Id:    "2",
			Name:  "Isadora",
			Email: "isadora@gmail.com",
		},
		{
			Id:    "3",
			Name:  "Ricardo",
			Email: "ricardo@gmail.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error on create stream: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 2)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}
	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) { // Sending and receiving with stream

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Could not make request: %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "Jeferson",
			Email: "jeferson@gmail.com",
		},
		{
			Id:    "2",
			Name:  "Isadora",
			Email: "isadora@gmail.com",
		},
		{
			Id:    "3",
			Name:  "Ricardo",
			Email: "ricardo@gmail.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user:", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error on receiving: %v", err)
			}

			fmt.Printf("Recebendo user %s com status: %s\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
