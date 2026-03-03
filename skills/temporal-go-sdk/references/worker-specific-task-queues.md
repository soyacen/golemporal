# Worker Specific Task Queues 特定 Worker 任务队列

为特定 Worker 分配任务队列。

```go
// Worker 1: 处理快速任务
func RunFastWorker() {
    client, _ := temporal.NewClient(...)

    worker := worker.New(client, "fast-task-queue", worker.Options{
        MaxConcurrentWorkflowExecution: 100,
        MaxConcurrentActivityExecution: 50,
    })

    worker.RegisterWorkflow(FastWorkflow)
    worker.RegisterActivity(FastActivity)

    worker.Run()
}

// Worker 2: 处理慢速任务
func RunSlowWorker() {
    client, _ := temporal.NewClient(...)

    worker := worker.New(client, "slow-task-queue", worker.Options{
        MaxConcurrentWorkflowExecution: 10,
        MaxConcurrentActivityExecution: 5,
    })

    worker.RegisterWorkflow(SlowWorkflow)
    worker.RegisterActivity(SlowActivity)

    worker.Run()
}

// 启动时指定任务队列
func StartWorkflowOnSpecificQueue(ctx context.Context, client temporal.Client) error {
    options := client.StartWorkflowOptions{
        TaskQueue: "fast-task-queue", // 指定任务队列
    }

    _, err := client.StartWorkflow(ctx, options, FastWorkflow, input)
    return err
}
```
