# Start Delay 延迟启动

延迟启动工作流。

```go
// 方式1: 使用 StartWorkflowOptions 的 StartDelay
func StartDelayedWorkflow(ctx context.Context, client temporal.Client) error {
    options := client.StartWorkflowOptions{
        ID:        "delayed-workflow-" + uuid.NewString(),
        TaskQueue: "my-task-queue",
        // 延迟启动
        StartDelay: time.Hour,
    }

    _, err := client.StartWorkflow(ctx, options, MyWorkflow, input)
    return err
}

// 方式2: 在工作流内部使用 Timer
func DelayedStartWorkflow(ctx workflow.Context, delay time.Duration) error {
    // 等待指定时间
    workflow.Sleep(ctx, delay)

    // 执行逻辑
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    return workflow.ExecuteActivity(ctx, MyActivity, nil).Get(ctx, nil)
}

// 方式3: 使用 Schedule 延迟
func ScheduleDelayedWorkflow(ctx context.Context, client temporal.Client) error {
    schedule := &schedule.Schedule{
        Spec: &schedule.Spec{
            // 延迟调度
            Delay: &durationpb.Duration{Seconds: 3600},
        },
        Action: &schedule.ScheduleAction{
            StartWorkflow: &schedule.StartWorkflow{
                Workflow: MyWorkflow,
                TaskQueue: "my-task-queue",
            },
        },
    }

    _, err := client.CreateSchedule(ctx, "delayed-schedule", schedule)
    return err
}
```
