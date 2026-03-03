# Cron Schedule 定时调度

使用 Cron 表达式创建周期性执行的工作流。

## 启动 Cron 工作流

```go
workflowOptions := client.StartWorkflowOptions{
    ID:           "cron-workflow-id",
    TaskQueue:    "cron-task-queue",
    CronSchedule: "@daily", // 每一天
}

client.ExecuteWorkflow(context.Background(), workflowOptions, CronWorkflow)
```

## Cron 表达式

| 表达式 | 说明 |
|--------|------|
| `@every 1m` | 每分钟 |
| `@every 5m` | 每5分钟 |
| `@hourly` | 每小时 |
| `@daily` | 每天午夜 |
| `@weekly` | 每周日 |
| `@monthly` | 每月1日 |

## 带参数的工作流

```go
// 每次执行时传入参数
workflowOptions := client.StartWorkflowOptions{
    CronSchedule: "@daily",
    Memo: map[string]interface{}{
        "lastRun": time.Now().Format(time.RFC3339),
    },
}
```

## Cron 历史

Temporal 会记录每次 Cron 执行的历史，便于审计。

## 最佳实践

1. **幂等性**: Cron 工作流必须幂等，因为可能重复执行
2. **状态检查**: 工作流开始时检查是否需要执行
3. **避免重叠**: 使用 Mutex 或 State 防止重叠执行
