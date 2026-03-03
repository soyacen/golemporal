# Request Response Query 请求响应查询

通过 Query 实现请求-响应模式。

```go
// 工作流状态
type RequestResponseState struct {
    Request  string
    Response string
    Done     bool
}

func RequestResponseQueryWorkflow(ctx workflow.Context) (*RequestResponseState, error) {
    state := &RequestResponseState{}

    // 注册 Query handler
    workflow.SetQueryHandler(ctx, "getState", func() (*RequestResponseState, error) {
        return state, nil
    })

    // 注册 Signal handler
    workflow.SetSignalHandler(ctx, "submit", func(data string) {
        state.Request = data
        state.Response = "processed: " + data
        state.Done = true
    })

    // 等待完成
    workflow.Await(ctx, func() bool { return state.Done })

    return state, nil
}

// 通过 Client Query 获取状态
func QueryWorkflowState(ctx context.Context, client temporal.Client, workflowID string) (*RequestResponseState, error) {
    resp, err := client.QueryWorkflow(ctx, workflowID, "", "getState")
    if err != nil {
        return nil, err
    }

    var state RequestResponseState
    err = resp.Get(&state)
    return &state, err
}
```
