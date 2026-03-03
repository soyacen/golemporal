# OpenTelemetry 观测性

集成 OpenTelemetry 进行分布式追踪。

## 安装依赖

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/jaeger
go get go.temporal.io/sdk/contrib/telemetry
```

## 配置 Tracing

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.temporal.io/sdk/contrib/telemetry"
)

exporter, err := jaeger.New(jaeger.WithCollectorEndpoint("http://jaeger:14268/api/traces"))
if err != nil {
    log.Fatal(err)
}

tp := trace.NewTracerProvider(trace.WithBatcher(exporter))
otel.SetTracerProvider(tp)

tracer := otel.Tracer("my-workflow")

// 配置 Client
c, err := client.Dial(client.Options{
    TracingInterceptorOptions: telemetry.TracingInterceptorOptions{
        Tracer: tracer,
    },
})
```

## Activity 追踪

```go
func TracedActivity(ctx context.Context, input string) (string, error) {
    tracer := otel.Tracer("my-activity")
    ctx, span := tracer.Start(ctx, "TracedActivity")
    defer span.End()

    // 业务逻辑
    return "result", nil
}
```

## Metrics

```go
meter := otel.Meter("my-app")

counter, _ := meter.Int64Counter("activity_executions")
counter.Add(ctx, 1)
```

## 最佳实践

1. **采样策略**: 配置合适的采样率
2. **Span 命名**: 使用描述性的 span 名称
3. **属性添加**: 添加业务相关属性便于调试
