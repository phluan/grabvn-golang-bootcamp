package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grab.com/luanpham/users_feedback/pb"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUsersClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Create(ctx, &pb.UserSignUpRequest{Email: "phluan92@gmail.com", Password: "123", FirstName: "Luan", LastName: "Pham"})
	if err != nil {
		log.Fatalf("could not create: %v", err)
	}
	log.Printf("Greeting: User#%s", r.Id)
}
