# External Environment Configuration 外部环境配置

从外部系统获取配置。

```go
// 配置 Activity
type Config struct {
    DatabaseURL string
    APIKey      string
    MaxRetries  int
}

func FetchConfigActivity(ctx context.Context, env string) (*Config, error) {
    // 从环境变量、配置服务或 Vault 获取配置
    config := &Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        APIKey:      os.Getenv("API_KEY"),
        MaxRetries:  3,
    }
    return config, nil
}

// 使用配置的工作流
func ConfiguredWorkflow(ctx workflow.Context, input Input) error {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    var config Config
    err := workflow.ExecuteActivity(ctx, FetchConfigActivity, "production").Get(ctx, &config)
    if err != nil {
        return err
    }

    // 使用配置执行后续逻辑
    return workflow.ExecuteActivity(ctx, ProcessWithConfig, config, input).Get(ctx, nil)
}
```
