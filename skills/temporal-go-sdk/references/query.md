# Query 查询工作流状态

使用 Query API 获取工作流的当前状态。

## 工作流端注册 Query

```go
func QueryableWorkflow(ctx workflow.Context) (int, error) {
    counter := 0

    // 注册查询处理程序
    workflow.SetQueryHandler(ctx, "get-counter", func() (int, error) {
        return counter, nil
    })

    // 另一个可参数化的查询
    workflow.SetQueryHandler(ctx, "get-history", func(n int) ([]string, error) {
        if n > 10 {
            n = 10
        }
        return history[:n], nil
    })

    // 等待信号
    signalChan := workflow.GetSignalChannel(ctx, "increment")
    for {
        signalChan.Receive(ctx, nil)
        counter++
    }
}
```

## 客户端查询

```go
// 同步查询
value, err := c.QueryWorkflow(context.Background(), "workflow-id", "", "get-counter")
if err != nil {
    return err
}
var count int
value.Get(&count)

// 带参数的查询
value, err := c.QueryWorkflow(context.Background(), "workflow-id", "", "get-history", 5)
var history []string
value.Get(&history)
```

## 查询工作流执行状态

```go
// 获取工作流元信息
we, err := c.DescribeWorkflowExecution(context.Background(), "workflow-id", "")
if err != nil {
    return err
}

fmt.Println("Status:", we.Status)
fmt.Println("StartTime:", we.StartTime)
fmt.Println("CloseTime:", we.CloseTime)
```

## 最佳实践

1. **只读操作**: Query 不应修改工作流状态
2. **幂等性**: Query 应该幂等
3. **快速响应**: Query 处理应该快速完成
4. **可选参数**: 支持可选参数增加灵活性
