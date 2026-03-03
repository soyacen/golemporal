# Standalone Activity 独立 Activity

独立运行的 Activity，不属于任何工作流。

```go
// Activity 定义（可以独立运行）
func StandaloneActivity(ctx context.Context, input Input) (Result, error) {
    // Activity 逻辑
    return Result{Value: "processed"}, nil
}

// 注册为独立 Activity
func init() {
    activity.Register(StandaloneActivity)
}

// 从外部调用独立 Activity
func CallStandaloneActivity(ctx context.Context, client temporal.Client) error {
    // 使用 WorkflowClient 调用 Activity
    workflowClient := client.WorkflowServiceClient()

    // 直接通过 workflowClient 调用 Activity（不启动工作流）
    // 注意: 这需要特殊的处理，通常通过启动一个"虚拟"工作流来实现
    _, err := workflowClient.StartWorkflow(ctx, &workflowservice.StartWorkflowExecutionRequest{
        WorkflowType: &workflowType{Name: "StandaloneActivityWrapper"},
        TaskQueue:   &taskQueue{Name: "standalone-activities"},
        Input:       marshal(input),
    })

    return err
}

// 或者使用 Local Activity
func WorkflowWithLocalActivity(ctx workflow.Context, input Input) error {
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute,
        // 标记为 Local Activity
        LocalProviderOptions: &workflow.LocalProviderOptions{
            ForceLocal: true,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    // Local Activity 在本地执行，不经过 Temporal Server
    return workflow.ExecuteLocalActivity(ctx, StandaloneActivity, input).Get(ctx, nil)
}
```
