# Child Workflow

子工作流允许将复杂业务流程分解为可管理的模块。

## 基本子工作流调用

```go
func ParentWorkflow(ctx workflow.Context, input ParentInput) (string, error) {
    // 设置子工作流选项
    cwo := workflow.ChildWorkflowOptions{
        WorkflowID:        "child-workflow-" + input.ID,
        TaskQueue:        "child-task-queue",
        ExecutionTimeout:  10 * time.Minute,
        RunTimeout:       5 * time.Minute,
    }
    ctx = workflow.WithChildOptions(ctx, cwo)

    var childResult string
    // 启动子工作流
    future := workflow.ExecuteChildWorkflow(ctx, ChildWorkflow, ChildInput{
        Data: input.Data,
    })

    // 等待子工作流完成
    err := future.Get(ctx, &childResult)
    if err != nil {
        return "", err
    }

    return "Parent completed with child result: " + childResult, nil
}
```

## ChildWorkflowOptions 配置

| 选项 | 说明 |
|------|------|
| WorkflowID | 子工作流 ID，可用于查询 |
| TaskQueue | 任务队列，默认使用父工作流的队列 |
| ExecutionTimeout | 执行总超时时间 |
| RunTimeout | 单次运行超时 |
| ParentClosePolicy | 父工作流关闭时子工作流的策略 |

## ParentClosePolicy

```go
cwo := workflow.ChildWorkflowOptions{
    // RequestCancel: 请求取消子工作流
    // Terminate: 终止子工作流
    // Abandon: 保持子工作流运行
    ParentClosePolicy: workflow.ParentClosePolicyRequestCancel,
}
```

## ContinueAsNew

用于处理大数据量，将状态转移到新的工作流执行：

```go
func LargeDataWorkflow(ctx workflow.Context, input LargeInput) (string, error) {
    // 处理第一批数据
    processBatch(input.Batch1)

    // 转移到新的工作流处理下一批
    workflow.ContinueAsNew(ctx, LargeInput{
        Batch1: input.Batch2,
        Batch2: input.Batch3,
    })
    return "", nil  // 不会执行
}
```

## 最佳实践

1. **模块化设计**: 子工作流应该是独立可重用的业务逻辑单元
2. **超时设置**: 总是设置合理的超时时间
3. **错误处理**: 处理子工作流可能的失败
4. **ID 管理**: 使用有意义的子工作流 ID 便于追踪
