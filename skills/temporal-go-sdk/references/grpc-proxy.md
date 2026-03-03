# gRPC Proxy gRPC 代理

通过 gRPC 代理与 Temporal Server 通信。

```go
// 创建 gRPC 拨号选项
func CreateGRPCProxyDialOptions(proxyAddr string) []grpc.DialOption {
    options := []grpc.DialOption{
        grpc.WithInsecure(),
        // 添加代理拦截器
        grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{},
            cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
            // 代理转发逻辑
            return invoker(ctx, method, req, reply, cc, opts...)
        }),
    }
    return options
}

// 使用代理连接
func CreateProxyClient(ctx context.Context, proxyAddr string) (temporal.Client, error) {
    dialOpts := CreateGRPCProxyDialOptions(proxyAddr)

    client, err := temporal.NewClient(temporal.ClientOptions{
        HostPort:  proxyAddr,
        DialOptions: dialOpts,
    })

    return client, err
}
```
