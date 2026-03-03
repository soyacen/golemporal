# Mutex 互斥锁

工作流中的互斥访问模式。

```go
func MutexWorkflow(ctx workflow.Context, resourceID string) error {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 使用 Semaphore 实现互斥
    semaphore := workflow.NewSemaphore(ctx, 1)

    // 获取锁
    err := semaphore.Acquire(ctx, func() error {
        // 临界区：一次只有一个工作流可以执行这段代码
        var result string
        return workflow.ExecuteActivity(ctx, AccessResourceActivity, resourceID).Get(ctx, &result)
    })

    return err
}

// 或者使用 Selector 实现
func MutexWithSelectorWorkflow(ctx workflow.Context, resourceID string) error {
    mutex := workflow.NewMutex(ctx)

    // 等待获取锁
    err := mutex.Lock(ctx, func() error {
        // 临界区
        var result string
        return workflow.ExecuteActivity(ctx, AccessResourceActivity, resourceID).Get(ctx, &result)
    })

    return err
}
```
