# Goroutine 并发执行

使用 `workflow.Go` 在工作流中并发执行多个任务。

## 并行执行多个 Activity

```go
func ParallelWorkflow(ctx workflow.Context, items []string) ([]string, error) {
    // 创建 Future 切片
    futures := make([]workflow.Future, len(items))

    // 并行启动所有 Activity
    for i, item := range items {
        futures[i] = workflow.ExecuteActivity(ctx, ProcessItem, item)
    }

    // 收集结果
    results := make([]string, len(items))
    for i, f := range futures {
        var result string
        if err := f.Get(ctx, &result); err != nil {
            return nil, err
        }
        results[i] = result
    }

    return results, nil
}
```

## 并行执行带错误处理

```go
func ParallelWithErrorHandling(ctx workflow.Context, items []string) (map[string]error, error) {
    futures := make([]workflow.Future, len(items))

    // 启动所有 Activity
    for i, item := range items {
        futures[i] = workflow.ExecuteActivity(ctx, ProcessItem, item)
    }

    // 收集结果和错误
    results := make(map[string]error)
    for i, item := range items {
        var result string
        err := futures[i].Get(ctx, &result)
        if err != nil {
            results[item] = err
        } else {
            results[item] = nil
        }
    }

    return results, nil
}
```

## 使用 Go Routine 执行副作用操作

```go
func WorkflowWithSideEffect(ctx workflow.Context) error {
    // workflow.Go 用于执行不依赖返回值的并发操作
    workflow.Go(ctx, func(ctx workflow.Context) {
        // 这个 goroutine 在后台运行
        // 注意：不能直接获取返回值
        logger := workflow.GetLogger(ctx)
        logger.Info("Background task running")
    })

    // 主流程继续执行
    return nil
}
```

## 扇出模式

```go
func FanOutWorkflow(ctx workflow.Context, batch Batch) error {
    // 将任务分成多个批次并行处理
    batchSize := 10
    var futures []workflow.Future

    for i := 0; i < len(batch.Items); i += batchSize {
        end := i + batchSize
        if end > len(batch.Items) {
            end = len(batch.Items)
        }
        subBatch := batch.Items[i:end]
        futures = append(futures, workflow.ExecuteActivity(ctx, ProcessBatch, subBatch))
    }

    // 等待所有批次完成
    for _, f := range futures {
        if err := f.Get(ctx, nil); err != nil {
            return err
        }
    }

    return nil
}
```

## 最佳实践

1. **控制并发数量**: 使用分批或信号量模式限制并发数
2. **错误处理**: 每个 Future 单独处理错误
3. **超时设置**: 为每个 Activity 设置合理的超时
4. **结果合并**: 使用 Future 收集结果后合并
