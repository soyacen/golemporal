# Greetings 问候示例

基本的问候工作流。

```go
// 简单的问候 Activity
func GreetActivity(ctx context.Context, name string) (string, error) {
    return "Hello, " + name + "!", nil
}

// 问候工作流
func GreetingWorkflow(ctx workflow.Context, name string) (string, error) {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Second * 30}
    ctx = workflow.WithActivityOptions(ctx, AO)

    var greeting string
    err := workflow.ExecuteActivity(ctx, GreetActivity, name).Get(ctx, &greeting)
    return greeting, err
}
```
