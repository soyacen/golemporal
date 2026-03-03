# Early Return 提前返回

工作流中的提前返回模式。

```go
func EarlyReturnWorkflow(ctx workflow.Context, input Input) (string, error) {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 快速路径检查
    if input.SkipProcessing {
        return "skipped", nil
    }

    // 前置验证
    if !input.IsValid() {
        return "", errors.New("invalid input")
    }

    // 主要业务逻辑
    var result string
    err := workflow.ExecuteActivity(ctx, ProcessActivity, input).Get(ctx, &result)
    if err != nil {
        return "", err
    }

    // 后置检查
    if result == "failed" {
        return "failed", errors.New("process failed")
    }

    return result, nil
}
```
