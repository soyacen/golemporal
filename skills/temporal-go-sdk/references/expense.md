# Expense 异步 Activity 完成

异步完成 Activity，用于长时间运行的任务。

## Activity 实现

```go
func ProcessExpenseActivity(ctx context.Context, expense Expense) error {
    // 模拟长时间处理
    activity.RecordHeartbeat(ctx, "processing")

    // 异步完成 - 立即返回，不阻塞
    // 稍后通过 client.CompleteActivity 异步完成
    return nil
}

// 异步完成的调用方
func AsyncCompleteWorkflow(ctx workflow.Context, expense Expense) error {
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: 1 * time.Minute,
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 不等待结果，立即返回
    workflow.ExecuteActivity(ctx, ProcessExpenseActivity, expense)

    // 继续执行其他逻辑
    return nil
}
```
