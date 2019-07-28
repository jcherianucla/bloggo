package main

import (
	"github.com/jcherianucla/bloggo/idl/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := proto.NewBloggoClient(conn)
	response, err := c.Create(context.Background(), &proto.CreatePostRequest{
		Title:       "Hello World",
		Description: "This is the first test",
	})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %v", response.Data)
}
