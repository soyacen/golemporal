# Synchronous Proxy 同步代理

同步代理模式。

```go
// 同步代理工作流
func SynchronousProxyWorkflow(ctx workflow.Context, request ProxyRequest) (ProxyResponse, error) {
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute,
        // 同步执行等待
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 调用远程服务
    var response ProxyResponse
    err := workflow.ExecuteActivity(ctx, ProxyActivity, request).Get(ctx, &response)
    if err != nil {
        return ProxyResponse{}, err
    }

    return response, nil
}

// 代理 Activity
func ProxyActivity(ctx context.Context, request ProxyRequest) (ProxyResponse, error) {
    // 调用实际的远程服务
    client := &http.Client{}
    req, _ := http.NewRequest(request.Method, request.URL, bytes.NewReader(request.Body))

    resp, err := client.Do(req)
    if err != nil {
        return ProxyResponse{}, err
    }

    body, _ := io.ReadAll(resp.Body)

    return ProxyResponse{
        StatusCode: resp.StatusCode,
        Body:       body,
    }, nil
}
```
