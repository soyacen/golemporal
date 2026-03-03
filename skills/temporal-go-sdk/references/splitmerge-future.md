# Split Merge Future 分离合并 Future

使用 Future 分离和合并任务。

```go
func SplitMergeFutureWorkflow(ctx workflow.Context, items []Item) ([]Result, error) {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 分离: 启动多个 Future
    futures := make([]workflow.Future, len(items))
    for i, item := range items {
        futures[i] = workflow.ExecuteActivity(ctx, ProcessItemActivity, item)
    }

    // 合并: 等待所有 Future 完成
    results := make([]Result, len(items))
    for i, f := range futures {
        err := f.Get(ctx, &results[i])
        if err != nil {
            return nil, err
        }
    }

    return results, nil
}
```
