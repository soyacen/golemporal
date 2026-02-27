package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/soyacen/golemporal/example/api"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// Create worker
	taskQueue := "golemporal-example"
	w := worker.New(c, taskQueue, worker.Options{})

	api.RegisterGreeterWorkflowWorker(w,
		&GreeterWorkflowServer{
			addActivity:   api.NewAddActivityClient(),
			multiActivity: api.NewMultiActivityClient(),
		},
		&AddActivityServer{},
		&MultiActivityServer{},
	)
	if err := w.Start(); err != nil {
		log.Fatalln("Unable to start worker", err)
	}

	// Wait for interrupt
	<-make(chan struct{})
}

type GreeterWorkflowServer struct {
	addActivity   api.AddActivityClient
	multiActivity api.MultiActivityClient
}

func (s *GreeterWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorkflow starting")
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	helloResult, err := s.addActivity.Add(ctx, &api.AddRequest{Count: input.GetCount()})
	if err != nil {
		logger.Error("activity failed", "error", err)
		return nil, err
	}
	logger.Info("HelloWorkflow completed")
	return &api.HelloResponse{
		Message: fmt.Sprintf("Hello, %s! (result: %d)", input.Name, helloResult.Result),
		Result:  helloResult.Result,
	}, nil
}

func (s *GreeterWorkflowServer) Goodbye(ctx workflow.Context, input *api.GoodbyeRequest) (*api.GoodbyeResponse, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("GoodbyeWorkflow starting")
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	helloResult, err := s.multiActivity.Multi(ctx, &api.MultiRequest{Count: input.GetCount()})
	if err != nil {
		logger.Error("activity failed", "error", err)
		return nil, err
	}
	logger.Info("GoodbyeWorkflow completed")
	return &api.GoodbyeResponse{
		Message: fmt.Sprintf("Goodbye, %s! (result: %d)", input.Name, helloResult.Result),
		Result:  helloResult.Result,
	}, nil
}

// Activity implementations
type AddActivityServer struct{}

func (s *AddActivityServer) Add(ctx context.Context, input *api.AddRequest) (*api.AddResponse, error) {
	return &api.AddResponse{Result: input.Count + input.Count}, nil
}

type MultiActivityServer struct{}

func (s *MultiActivityServer) Multi(ctx context.Context, input *api.MultiRequest) (*api.MultiResponse, error) {
	return &api.MultiResponse{Result: input.Count * input.Count}, nil
}
