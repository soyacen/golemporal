# Memo 工作流元数据

在工作流中使用 Memo 存储额外信息。

```go
// 启动带 Memo 的工作流
func StartWorkflowWithMemo(ctx context.Context, client temporal.Client) error {
    memo := map[string]interface{}{
        "customer_id": "customer-123",
        "order_id":    "order-456",
        "priority":    "high",
    }

    options := client.StartWorkflowOptions{
        ID:        "workflow-" + uuid.NewString(),
        TaskQueue: "my-task-queue",
        Memo:      memo,
    }

    _, err := client.StartWorkflow(ctx, options, MyWorkflow, input)
    return err
}

// 在工作流中获取 Memo
func WorkflowWithMemo(ctx workflow.Context) error {
    // 通过 info 获取
    info := workflow.GetInfo(ctx)
    memo := info.Memo

    customerID := memo["customer_id"]
    // 使用 memo 数据
    return nil
}
```
