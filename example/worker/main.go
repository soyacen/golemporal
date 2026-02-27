package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/soyacen/golemporal/example"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
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

	ws := GreeterWorkflowServer{
		activity: example.NewGreeterActivityClient(),
	}
	example.RegisterGreeterWorker(w, &ws)

	as := GreeterActivityServer{}
	example.RegisterGreeterActivity(w, &as)

	if err := w.Start(); err != nil {
		log.Fatalln("Unable to start worker", err)
	}

	// Wait for interrupt
	<-make(chan struct{})
}

type GreeterWorkflowServer struct {
	activity example.GreeterActivityClient
}

func (s *GreeterWorkflowServer) HelloWorkflow(ctx workflow.Context, input *example.HelloRequest) (*example.HelloResponse, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("starting HelloWorkflow", zap.String("name", input.Name))
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	helloResult, err := s.activity.SayHello(ctx, input)
	if err != nil {
		logger.Error("activity failed", zap.Error(err))
		return nil, err
	}
	logger.Info("workflow completed", zap.String("message", helloResult.Message))
	return helloResult, nil
}

// Activity implementations
type GreeterActivityServer struct{}

func (s *GreeterActivityServer) SayHello(ctx context.Context, input *example.HelloRequest) (*example.HelloResponse, error) {
	return &example.HelloResponse{
		Message: fmt.Sprintf("Hello, %s! (count: %d)", input.Name, input.Count),
		Result:  input.Count * 2,
	}, nil
}
