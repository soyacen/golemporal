# Hello World 示例

基础工作流和 Activity 定义。

## 工作流定义

```go
package helloworld

import (
    "fmt"
    "time"

    "go.temporal.io/sdk/workflow"
)

func HelloWorldWorkflow(ctx workflow.Context, name string) (string, error) {
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Second,
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    var result string
    err := workflow.ExecuteActivity(ctx, HelloWorldActivity, name).Get(ctx, &result)
    return result, err
}
```

## Activity 定义

```go
func HelloWorldActivity(ctx context.Context, name string) (string, error) {
    return fmt.Sprintf("Hello, %s!", name), nil
}
```

## Worker 注册

```go
func main() {
    c, err := client.Dial(client.Options{})
    if err != nil {
        log.Fatalln("Unable to create client", err)
    }
    defer c.Close()

    w := worker.New(c, "hello-world-task-queue", worker.Options{})
    w.RegisterWorkflow(HelloWorldWorkflow)
    w.RegisterActivity(HelloWorldActivity)

    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalln("Worker stopped", err)
    }
}
```

## 启动工作流

```go
workflowOptions := client.StartWorkflowOptions{
    ID:        "hello-world-workflow-id",
    TaskQueue: "hello-world-task-queue",
}

result, err := c.ExecuteWorkflow(
    context.Background(),
    workflowOptions,
    HelloWorldWorkflow,
    "World",
)
```

## 关键点

1. **Workflow 函数签名**: `func(ctx workflow.Context, input InputType) (OutputType, error)`
2. **Activity 函数签名**: `func(ctx context.Context, input InputType) (OutputType, error)`
3. **ActivityOptions**: 必须设置超时时间
4. **ExecuteActivity**: 异步执行 Activity，通过 Future 获取结果
