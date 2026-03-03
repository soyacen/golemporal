# Request Response Update 请求响应更新

通过 Update 实现请求-响应模式。

```go
// 工作流
func RequestResponseUpdateWorkflow(ctx workflow.Context) (string, error) {
    result := ""

    // 注册 Update handler
    handler := workflow.UpdateWorkflow(ctx, "process", func(input string) string {
        // 处理请求
        result = "processed: " + input
        return result
    })

    // 等待 Update 触发
    workflow.Await(ctx, func() bool { return result != "" })

    return result, nil
}

// 客户端调用 Update
func CallUpdate(ctx context.Context, client temporal.Client) (string, error) {
    handle, err := client.UpdateWorkflow(ctx, "workflow-id", "", "process", "my request")
    if err != nil {
        return "", err
    }

    var response string
    err = handle.Get(ctx, &response)
    return response, err
}
```
