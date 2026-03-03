# Split Merge Selector 分离合并选择器

使用 Selector 模式分离和合并任务。

```go
func SplitMergeSelectorWorkflow(ctx workflow.Context, items []Item) ([]Result, error) {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 分离: 启动多个 Activity
    futures := make([]workflow.Future, len(items))
    for i, item := range items {
        futures[i] = workflow.ExecuteActivity(ctx, ProcessItemActivity, item)
    }

    // 使用 Selector 等待所有完成
    results := make([]Result, len(items))
    completed := 0

    selector := workflow.NewSelector(ctx)
    for i, f := range futures {
        idx := i
        future := f
        selector.AddFuture(f, func(f workflow.Future) {
            err := f.Get(ctx, &results[idx])
            if err != nil {
                // 处理错误
            }
            completed++
        })
    }

    // 等待所有完成
    for completed < len(items) {
        selector.Select(ctx)
    }

    return results, nil
}
```
