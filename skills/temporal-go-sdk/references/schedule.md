# Schedule 定时调度

使用 Schedule 定时执行工作流。

```go
// 创建定时调度
func CreateSchedule(ctx context.Context, client temporal.Client) error {
    schedule := &schedule.Schedule{
        Spec: &schedule.Spec{
            // 每天上午 9 点执行
            Calendars: []schedule.CalendarSpec{
                {
                    Hour:    []int{9},
                    Minute:  []int{0},
                    Weekday: []schedule.Weekday{schedule.Monday, schedule.Tuesday, schedule.Wednesday, schedule.Thursday, schedule.Friday},
                },
            },
        },
        Action: &schedule.ScheduleAction{
            StartWorkflow: &schedule.StartWorkflow{
                ID:        "daily-report-workflow",
                TaskQueue:  "reports",
                Workflow:   DailyReportWorkflow,
                Input:     nil,
            },
        },
    }

    handle, err := client.CreateSchedule(ctx, "my-schedule", schedule)
    return err
}

// 调度工作流
func DailyReportWorkflow(ctx workflow.Context) error {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Hour}
    ctx = workflow.WithActivityOptions(ctx, AO)

    return workflow.ExecuteActivity(ctx, GenerateReportActivity).Get(ctx, nil)
}
```
