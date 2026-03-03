# Zap Adapter Zap 日志适配器

使用 Zap 日志库。

```go
import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// 创建 Zap logger
func NewZapLogger() *zap.Logger {
    config := zap.Config{
        Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
        Development: false,
        Encoding:    "json",
        EncoderConfig: zapcore.EncoderConfig{
            TimeKey:        "ts",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "caller",
            MessageKey:     "msg",
            StacktraceKey:  "stacktrace",
            LineEnding:     zapcore.DefaultLineEnding,
            EncodeLevel:    zapcore.LowercaseLevelEncoder,
            EncodeTime:     zapcore.ISO8601TimeEncoder,
            EncodeDuration: zapcore.SecondsDurationEncoder,
            EncodeCaller:   zapcore.ShortCallerEncoder,
        },
        OutputPaths:      []string{"stdout"},
        ErrorOutputPaths: []string{"stderr"},
    }

    logger, _ := config.Build()
    return logger
}

// 使用 Zap 的 Worker
func RunWorkerWithZap(ctx context.Context) error {
    logger := NewZapLogger()
    defer logger.Sync()

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
```
