# Choice Multi 多条件分支

多个条件分支的并行或顺序执行。

```go
func ChoiceMultiWorkflow(ctx workflow.Context, input Input) (string, error) {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 并行执行多个条件分支
    futures := []workflow.Future{}

    if input.CheckA {
        f := workflow.ExecuteActivity(ctx, ActivityA, input.Data)
        futures = append(futures, f)
    }
    if input.CheckB {
        f := workflow.ExecuteActivity(ctx, ActivityB, input.Data)
        futures = append(futures, f)
    }
    if input.CheckC {
        f := workflow.ExecuteActivity(ctx, ActivityC, input.Data)
        futures = append(futures, f)
    }

    // 等待所有分支完成
    for _, f := range futures {
        if err := f.Get(ctx, nil); err != nil {
            return "", err
        }
    }

    return "all completed", nil
}
```
