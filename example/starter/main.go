package main

import (
	"context"
	"log"

	"github.com/soyacen/golemporal/example/api"
	"go.temporal.io/sdk/client"
)

func main() {
	ctx := context.Background()
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	taskQueue := "golemporal-example"

	gc := api.NewGreeterWorkflowClient(c, taskQueue)
	helloResult, err := gc.Hello(ctx, &api.HelloRequest{Name: "World", Count: 5})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("hello workflow completed message: %s, result: %d", helloResult.Message, helloResult.Result)

	goodbyeResult, err := gc.Goodbye(ctx, &api.GoodbyeRequest{Name: "World", Count: 10})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("goodbye workflow completed message: %s, result: %d", goodbyeResult.Message, goodbyeResult.Result)
}
