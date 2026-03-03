# Datadog 集成

集成 Datadog 监控。

```go
import (
    "github.com/DataDog/datadog-go/statsd"
)

// 创建 Datadog 客户端
func NewDatadogClient() (*statsd.Client, error) {
    return statsd.New("localhost:8125", statsd.WithNamespace("temporal."))
}

// Datadog 拦截器
type DatadogInterceptor struct {
    statsd *statsd.Client
}

func NewDatadogInterceptor() (*DatadogInterceptor, error) {
    statsdClient, err := NewDatadogClient()
    if err != nil {
        return nil, err
    }

    return &DatadogInterceptor{statsd: statsdClient}, nil
}

func (d *DatadogInterceptor) InterceptActivity(ctx context.Context, method string,
    req interface{}, info *activity.Info, next activity.NextInterceptor) (interface{}, error) {

    // 记录开始时间
    start := time.Now()

    // 记录活动调用
    d.statsd.Incr("activity.started", []string{"method:" + method}, 1)

    result, err := next.Execute(ctx, req)

    // 记录结果
    tags := []string{"method:" + method}
    if err != nil {
        d.statsd.Incr("activity.failed", tags, 1)
    } else {
        d.statsd.Incr("activity.completed", tags, 1)
    }

    // 记录执行时间
    d.statsd.Timing("activity.duration", time.Since(start), tags, 1)

    return result, err
}

// 使用拦截器
func RunWorkerWithDatadog(ctx context.Context) error {
    interceptor, err := NewDatadogInterceptor()
    if err != nil {
        return err
    }

    client, _ := temporal.NewClient(...)

    worker := worker.New(client, "my-task-queue", worker.Options{
        Interceptors: []interceptor.Interceptor{interceptor},
    })

    return worker.Run()
}
```
