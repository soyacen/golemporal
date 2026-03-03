# Nexus 跨命名空间调用

使用 Nexus API 实现跨命名空间的工作流调用。

## 客户端配置 Nexus

```go
c, err := client.Dial(client.Options{
    HostPort:  "remote-cluster.tmprl.cloud:7233",
    Namespace: "my-namespace",
})
```

## 启动远程工作流

```go
// 通过 Nexus 启动另一个命名空间的工作流
handle, err := c.StartWorkflowNexus(context.Background(), client.StartWorkflowNexusOptions{
    Namespace: "other-namespace",
    Workflow: "RemoteWorkflow",
    TaskQueue: "remote-task-queue",
    Input:    myInput,
})
```

## Nexus 取消

参考: nexus-cancelation 示例

```go
// 取消 Nexus 操作
err := handle.Cancel(context.Background())
```

## Nexus 上下文传播

参考: nexus-context-propagation 示例

```go
// 传播上下文
ctx := context.WithValue(context.Background(), "key", "value")
handle, err := c.StartWorkflowNexus(ctx, client.StartWorkflowNexusOptions{
    // ...
})
```

## 最佳实践

1. **安全**: 使用 mTLS 保护跨集群通信
2. **超时**: 设置合理的超时时间
3. **错误处理**: 处理网络错误和超时
