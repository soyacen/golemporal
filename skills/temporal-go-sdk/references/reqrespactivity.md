# ReqResp Activity 请求响应模式

通过信号接收请求，通过 Activity 回调返回响应。

```go
// 工作流实现
type RequestResponseWorkflow struct {
    Result string
}

func (wf *RequestResponseWorkflow) Execute(ctx workflow.Context) (string, error) {
    // 注册信号处理程序接收请求
    workflow.SetSignalHandler(ctx, "request", func(data string) {
        wf.Result = data
    })

    // 等待信号
    workflow.Await(ctx, func() bool { return wf.Result != "" })
    return wf.Result, nil
}

// 请求 Activity - 接收请求
func RequestActivity(ctx context.Context, request Request) error {
    // 处理请求逻辑
    return nil
}

// 响应 Activity - 回调返回结果
func ResponseActivity(ctx context.Context, response Response) error {
    return nil
}
```
