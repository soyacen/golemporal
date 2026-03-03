# Slog Adapter Slog 适配器

使用 slog 日志。

```go
import (
    "context"
    "log/slog"
)

// 创建 slog logger
func NewSlogLogger() *slog.Logger {
    return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
}

// 使用 slog 的 Worker
func RunWorkerWithSlog(ctx context.Context) error {
    logger := NewSlogLogger()

    client, err := temporal.NewClient(temporal.ClientOptions{
        Logger: logger,
    })
    if err != nil {
        return err
    }

    worker := worker.New(client, "my-task-queue", worker.Options{
        Logger: logger,
    })

    worker.RegisterWorkflow(MyWorkflow)
    worker.RegisterActivity(MyActivity)

    return worker.Run(background.Background())
}

// Activity 中使用 slog
func ActivityWithSlog(ctx context.Context, input Input) error {
    logger := slog.Default()

    logger.Info("activity started", "input", input)

    // 处理逻辑

    logger.Info("activity completed")
    return nil
}
```
