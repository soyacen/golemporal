# Worker Versioning 工作流版本控制

使用 Worker Versioning 实现无停机部署。

## 配置 Versioning

```go
w := worker.New(c, "task-queue", worker.Options{
    BuildID:        "build-v1",
    UseBuildID:     true,
})
```

## 部署新版本

```go
// 新版本
w := worker.New(c, "task-queue", worker.Options{
    BuildID:        "build-v2",
    UseBuildID:     true,
    BuildIDMode:     worker.AutoUpgrade,
})
```

## BuildID 模式

| 模式 | 说明 |
|------|------|
| AutoUpgrade | 自动升级到最新版本 |
| Immutable | 固定版本，不升级 |
| PinLatest | 固定使用最新版本 |

## 兼容版本

```go
w := worker.New(c, "task-queue", worker.Options{
    BuildID:             "build-v2",
    UseBuildID:          true,
    BuildIDMode:         worker.AutoUpgrade,
    AllowedBuildIDs:     []string{"build-v1"}, // 允许回退到 v1
})
```

## 最佳实践

1. **分阶段部署**: 逐步升级而非一次性全部升级
2. **回退支持**: 保留旧版本一段时间以便回退
3. **版本标记**: 使用语义化版本号
