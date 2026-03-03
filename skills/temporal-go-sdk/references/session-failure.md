# Session Failure 会话失败处理

Session 失败时的处理。

```go
// 使用 Session 的工作流
func SessionWorkflow(ctx workflow.Context, input Input) error {
    sessionOpts := &workflow.SessionOptions{
        ExecutionTimeout: time.Hour,
        CreationTimeout: time.Minute * 10,
    }

    sessionCtx, err := workflow.CreateSession(ctx, sessionOpts)
    if err != nil {
        return err
    }

    defer func() {
        // 确保 Session 清理
        workflow.CompleteSession(sessionCtx)
    }()

    // 执行多个 Activity
    for _, item := range input.Items {
        err := workflow.ExecuteActivity(sessionCtx, ProcessItemActivity, item).Get(sessionCtx, nil)
        if err != nil {
            // Session 失败处理
            return handleSessionFailure(sessionCtx, err)
        }
    }

    return nil
}

func handleSessionFailure(sessionCtx workflow.Context, err error) error {
    // 检查是否是 Session 错误
    if sessionCtx.Err() != nil {
        // 重建 Session
        newSessionCtx, recreateErr := workflow.RecreateSession(sessionCtx, &workflow.SessionOptions{
            ExecutionTimeout: time.Hour,
            CreationTimeout: time.Minute * 10,
        })

        if recreateErr != nil {
            return recreateErr
        }

        // 在新 Session 中重试
        return retryInNewSession(newSessionCtx)
    }

    return err
}
```
