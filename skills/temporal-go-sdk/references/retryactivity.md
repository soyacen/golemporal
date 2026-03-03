# Activity 重试策略

使用 RetryPolicy 配置 Activity 的重试行为。

## 基本重试配置

```go
AO := workflow.ActivityOptions{
    StartToCloseTimeout: 5 * time.Minute,
    RetryPolicy: &temporal.RetryPolicy{
        InitialInterval:    time.Second,        // 初始间隔
        BackoffCoefficient: 2.0,                 // 退避系数
        MaximumInterval:    100 * time.Second,   // 最大间隔
        MaximumAttempts:   5,                   // 最大重试次数
    },
}
ctx = workflow.WithActivityOptions(ctx, AO)
```

## 带心跳的重试

长运行 Activity 应该使用心跳来跟踪进度：

```go
// Activity 实现
func LongRunningActivity(ctx context.Context, taskID string) error {
    for i := 0; i < 100; i++ {
        // 检查是否被取消
        if activity.IsActivity(ctx) {
            activity.RecordHeartbeat(ctx, i)  // 记录进度
        }
        time.Sleep(time.Second)
    }
    return nil
}

// 获取心跳详情用于恢复
func LongRunningActivityWithRecovery(ctx context.Context, taskID string) error {
    var progress int
    if activity.GetHeartbeatDetails(ctx, &progress) == nil {
        // 从上次进度恢复
        fmt.Printf("Resuming from progress: %d\n", progress)
    }
    // ... 继续处理
}
```

## 重试策略参数

| 参数 | 说明 | 示例 |
|------|------|------|
| InitialInterval | 第一次重试等待时间 | 1s |
| BackoffCoefficient | 间隔增长系数 | 2.0 |
| MaximumInterval | 最大等待间隔 | 100s |
| MaximumAttempts | 最大重试次数，0 表示无限 | 5 |
| NonRetryableErrorTypes | 不重试的错误类型 | []string{"InvalidArgument"} |

## 非重试错误

```go
RetryPolicy: &temporal.RetryPolicy{
    MaximumAttempts: 3,
    NonRetryableErrorTypes: []string{
        "InvalidArgument",
        "NotFound",
        "Unauthorized",
    },
}
```

## 最佳实践

1. **始终设置 MaximumAttempts**: 防止无限重试
2. **使用心跳**: 长运行 Activity 必须心跳
3. **幂等性**: Activity 设计要支持重试
4. **区分可重试和不可重试错误**: 使用 NonRetryableErrorTypes
