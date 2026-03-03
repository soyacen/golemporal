# Eager Workflow Start 积极启动工作流

立即在调用方 Worker 上启动工作流，而不是等待 Server 调度。

```go
// 使用 StartWorkflowOptions 启用积极启动
func StartEagerWorkflow(ctx context.Context, client temporal.Client) error {
    options := client.StartWorkflowOptions{
        ID:        "workflow-id-" + uuid.NewString(),
        TaskQueue: "your-task-queue",
        // 启用积极启动
        StartLimits: workflow.StartLimits{
            EagerStart: true,
        },
    }

    // 立即返回 WorkflowRun，不需要等待调度
    run, err := client.StartWorkflow(ctx, options, MyWorkflow, input)
    if err != nil {
        return err
    }

    // 可以立即获取结果（如果工作流很快完成）
    var result string
    err = run.Get(ctx, &result)
    return err
}
```
