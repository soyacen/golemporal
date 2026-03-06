package main

import (
	"context"
	"log"

	"github.com/soyacen/golemporal/example/api"
	"github.com/soyacen/golemporal/starter"
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

	hc := api.NewHelloWorkflowClient(c, taskQueue)
	var helloResult api.HelloResponse
	if err := hc.Hello(ctx, &api.HelloRequest{Name: "World", Count: 5}, starter.GetResult(&helloResult)); err != nil {
		log.Fatal(err)
	}

	log.Printf("hello workflow completed message: %s, result: %d", helloResult.Message, helloResult.Result)

	gc := api.NewGoodbyeWorkflowClient(c, taskQueue)
	var goodbyeResult api.GoodbyeResponse
	if err := gc.Goodbye(ctx, &api.GoodbyeRequest{Name: "World", Count: 10}, starter.GetResult(&goodbyeResult)); err != nil {
		log.Fatal(err)
	}
	log.Printf("goodbye workflow completed message: %s, result: %d", goodbyeResult.Message, goodbyeResult.Result)
}
