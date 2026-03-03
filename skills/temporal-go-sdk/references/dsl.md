# DSL Workflow DSL 工作流

使用 DSL 定义工作流逻辑。

```go
// DSL 定义示例
type DSLWorkflow struct {
    Steps []Step
}

type Step struct {
    Name   string
    Type   string // "activity", "condition", "parallel"
    Action string
    Next   string
}

// DSL 解析器
func ParseDSL(yamlContent string) (*DSLWorkflow, error) {
    // 解析 YAML 定义
    // 返回可执行的 DSLWorkflow
}

// DSL 执行器
func ExecuteDSLWorkflow(ctx workflow.Context, dsl *DSLWorkflow) error {
    currentStep := dsl.Steps[0]

    for currentStep != nil {
        switch currentStep.Type {
        case "activity":
            err := workflow.ExecuteActivity(ctx,
                getActivityByName(currentStep.Action),
                nil).Get(ctx, nil)
            if err != nil {
                return err
            }
        case "condition":
            // 条件判断逻辑
        }

        currentStep = getNextStep(dsl, currentStep.Next)
    }

    return nil
}
```
