# Nexus Multiple Arguments Nexus 多参数

Nexus 操作传递多个参数。

```go
// 复杂请求类型
type ComplexRequest struct {
    UserID    string
    Action    string
    Data      map[string]interface{}
    Metadata  map[string]string
}

// 多参数 Nexus 调用
func CallWithMultipleArgs(ctx context.Context, client temporal.Client) error {
    request := ComplexRequest{
        UserID:   "user-123",
        Action:   "process",
        Data: map[string]interface{}{
            "amount":   100.0,
            "currency": "USD",
        },
        Metadata: map[string]string{
            "source": "api",
        },
    }

    handle, err := client.NexusStartOperation(ctx, nexus.StartOperationOptions{
        Service:   "payment-service",
        Operation: "process-payment",
        Request:   request,
    })

    if err != nil {
        return err
    }

    var result PaymentResult
    return handle.Get(ctx, &result)
}
```
