# Choice Exclusive 条件分支

根据输入动态选择执行路径。

```go
func ChoiceWorkflow(ctx workflow.Context, input Input) (string, error) {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 动态选择 Activity
    switch input.Choice {
    case "A":
        return workflow.ExecuteActivity(ctx, ActivityA, input).Get(ctx, nil)
    case "B":
        return workflow.ExecuteActivity(ctx, ActivityB, input).Get(ctx, nil)
    case "C":
        return workflow.ExecuteActivity(ctx, ActivityC, input).Get(ctx, nil)
    default:
        return "default", nil
    }
}
```
