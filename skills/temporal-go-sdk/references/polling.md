# Polling 轮询

周期性轮询外部资源。

```go
func PollingWorkflow(ctx workflow.Context, target string) (string, error) {
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: 30 * time.Second,
        RetryPolicy: &temporal.RetryPolicy{
            MaximumAttempts: 10,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    var result string
    err := workflow.ExecuteActivity(ctx, PollExternalService, target).Get(ctx, &result)
    return result, err
}

func PollExternalService(ctx context.Context, target string) (string, error) {
    for i := 0; i < 10; i++ {
        // 检查外部资源状态
        status, err := checkStatus(target)
        if err == nil && status == "ready" {
            return status, nil
        }

        // 记录心跳
        activity.RecordHeartbeat(ctx, fmt.Sprintf("attempt-%d", i))

        // 等待后重试
        time.Sleep(5 * time.Second)
    }
    return "", errors.New("timeout waiting for resource")
}
```
