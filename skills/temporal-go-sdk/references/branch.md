# Branch 分支执行

根据条件执行不同的 Activity 分支。

```go
func BranchWorkflow(ctx workflow.Context, input Input) (string, error) {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    switch input.Type {
    case "typeA":
        var result string
        err := workflow.ExecuteActivity(ctx, ActivityA, input.Data).Get(ctx, &result)
        return result, err
    case "typeB":
        var result string
        err := workflow.ExecuteActivity(ctx, ActivityB, input.Data).Get(ctx, &result)
        return result, err
    default:
        var result string
        err := workflow.ExecuteActivity(ctx, ActivityDefault, input.Data).Get(ctx, &result)
        return result, err
    }
}
```
