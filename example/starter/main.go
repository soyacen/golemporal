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
	helloResult, hlMd, err := hc.Hello(ctx, &api.HelloRequest{Name: "World", Count: 5}, starter.WaitResult(true))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("hello workflow completed, workflow_id: %s, workflow_type: %s, run_id: %s, message: %s, result: %d", hlMd.GetWorkflowId(), hlMd.GetWorkflowType(), hlMd.GetRunId(), helloResult.Message, helloResult.Result)

	gc := api.NewGoodbyeWorkflowClient(c, taskQueue)
	goodbyeResult, bgMd, err := gc.Goodbye(ctx, &api.GoodbyeRequest{Name: "World", Count: 10}, starter.WaitResult(true))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("goodbye workflow completed , workflow_id: %s, workflow_type: %s, run_id: %s, message: %s, result: %d", bgMd.GetWorkflowId(), bgMd.GetWorkflowType(), bgMd.GetRunId(), goodbyeResult.Message, goodbyeResult.Result)
}
