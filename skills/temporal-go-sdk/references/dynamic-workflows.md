# Dynamic Workflow 动态工作流

运行时动态决定执行哪个工作流。

```go
// 动态工作流入口
func DynamicWorkflow(ctx workflow.Context, input map[string]interface{}) (interface{}, error) {
    // 从输入中获取工作流类型
    workflowType, ok := input["workflowType"].(string)
    if !ok {
        return nil, errors.New("workflowType not found")
    }

    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 动态执行不同的逻辑
    switch workflowType {
    case "typeA":
        var result string
        err := workflow.ExecuteActivity(ctx, ActivityA, input).Get(ctx, &result)
        return result, err
    case "typeB":
        var result int
        err := workflow.ExecuteActivity(ctx, ActivityB, input).Get(ctx, &result)
        return result, err
    default:
        return nil, fmt.Errorf("unknown workflow type: %s", workflowType)
    }
}

// 注册动态工作流
func init() {
    workflow.Register(DynamicWorkflow)
}
```
