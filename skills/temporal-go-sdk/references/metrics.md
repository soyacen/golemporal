# Metrics 指标监控

使用 Prometheus 和 Tally 收集 Temporal 指标。

## 安装依赖

```bash
go get github.com/uber-go/tally
go get github.com/prometheus/client_golang/prometheus
```

## 配置 Metrics

```go
import (
    "github.com/uber-go/tally"
    "github.com/uber-go/tally/prometheus"
    "go.temporal.io/sdk/metrics"
)

reporter := prometheus.NewReporter(prometheus.Configuration{
    Prefix: "temporal_",
})

scope, closer := tally.NewRootScope(tally.ScopeOptions{
    Reporter:  reporter,
    Prefix:    "my-worker",
}, time.Second)
defer closer.Close()

// 配置 Client
c, err := client.Dial(client.Options{
    MetricsHandler: metrics.NewTallyHandler(scope),
})

// 配置 Worker
w := worker.New(c, "task-queue", worker.Options{
    MetricsHandler: metrics.NewTallyHandler(scope),
})
```

## 常用指标

| 指标 | 类型 | 说明 |
|------|------|------|
| workflow_execution_start | Counter | 工作流启动数 |
| activity_execution_start | Counter | Activity 启动数 |
| workflow_task_queue_poll_succeed | Counter | 任务队列轮询成功 |
| worker_poll_succeed | Counter | Worker 轮询成功 |

## 自定义指标

```go
counter := scope.Counter("custom_counter")
counter.Inc(1)

gauge := scope.Gauge("custom_gauge")
gauge.Set(1.5)

timer := scope.Timer("custom_timer")
timer.Record(time.Second)
```

## 最佳实践

1. **指标聚合**: 使用 Tally 的 scope 进行指标聚合
2. **标签**: 使用 Tags 区分不同维度
3. **采样**: 对于高频指标考虑采样
