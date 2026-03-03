# Greetings Local 本地问候

本地 Activity 调用（同一 Worker）。

```go
// Activity 定义
func GreetLocalActivity(ctx context.Context, name string) (string, error) {
    return "Hello, " + name + "!", nil
}

// 工作流
func GreetingLocalWorkflow(ctx workflow.Context, name string) (string, error) {
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: time.Second * 30,
        // 本地调度：优先在同一 Worker 上执行
        LocalProviderOptions: &workflow.LocalProviderOptions{
            ForceLocal: true,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    var greeting string
    err := workflow.ExecuteActivity(ctx, GreetLocalActivity, name).Get(ctx, &greeting)
    return greeting, err
}
```
