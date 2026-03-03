# Session 资源池

使用 Session 将多个 Activity 绑定到同一个 Worker。

## 创建 Session

```go
sessionOptions := &workflow.SessionOptions{
    ExecutionTimeout: 10 * time.Minute,
    CreationTimeout:  10 * time.Minute,
}

session, err := workflow.CreateSession(ctx, sessionOptions)
if err != nil {
    return err
}
defer session.Close()

// 所有在此 session 中执行的 Activity 将在同一 Worker 上运行
err = workflow.ExecuteActivity(session, Activity1, input).Get(ctx, nil)
err = workflow.ExecuteActivity(session, Activity2, input).Get(ctx, nil)
```

## Session 失败处理

```go
session, err := workflow.CreateSession(ctx, sessionOptions)
if err != nil {
    return err
}

defer func() {
    if session != nil {
        session.Close()
    }
}()

// 如果 Activity 失败，session 会被标记为失败
// 可以检查并重新创建
if session.GetStatus() == workflow.SessionStatusFailed {
    session, err = workflow.CreateSession(ctx, sessionOptions)
}
```

## 并发 Session 限制

```go
// 使用 Semaphore 限制并发 session 数量
sessionLimit := 5
sem := workflow.NewSemaphore(ctx, sessionLimit)

if sem.Acquire(ctx) {
    session, err := workflow.CreateSession(ctx, sessionOptions)
    defer session.Close()
    // 执行 activities
    sem.Release()
}
```

## 最佳实践

1. **资源管理**: 使用 defer 确保 session 正确关闭
2. **错误处理**: 处理 session 失败并重试
3. **并发控制**: 使用 Semaphore 限制并发 session 数量
