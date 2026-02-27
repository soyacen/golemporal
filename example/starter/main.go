package main

import (
	"context"
	"log"

	"github.com/soyacen/golemporal/example/api"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	taskQueue := "golemporal-example"

	gc := api.NewGreeterWorkflowClient(c, taskQueue)
	helloResult, err := gc.Hello(ctx, &api.HelloRequest{Name: "World", Count: 5})
	if err != nil {
		logger.Fatal("failed to execute workflow", zap.Error(err))
	}

	logger.Info("hello workflow completed",
		zap.String("message", helloResult.Message),
		zap.Int32("result", helloResult.Result))

	goodbyeResult, err := gc.Goodbye(ctx, &api.GoodbyeRequest{Name: "World", Count: 10})
	if err != nil {
		logger.Fatal("failed to execute workflow", zap.Error(err))
	}

	logger.Info("goodbye workflow completed",
		zap.String("message", goodbyeResult.Message),
		zap.Int32("result", goodbyeResult.Result))
}
