# Dynamic Config 动态配置

动态配置管理。

```go
// 动态配置源
type DynamicConfigSource struct {
    config map[string]interface{}
    mutex  sync.RWMutex
}

func NewDynamicConfigSource() *DynamicConfigSource {
    return &DynamicConfigSource{
        config: map[string]interface{}{
            "max_retries":      3,
            "timeout_seconds":  60,
            "enable_feature_x": false,
        },
    }
}

// 获取配置
func (d *DynamicConfigSource) Get(key string) interface{} {
    d.mutex.RLock()
    defer d.mutex.RUnlock()
    return d.config[key]
}

// 更新配置
func (d *DynamicConfigSource) Update(key string, value interface{}) {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    d.config[key] = value
}

// 在工作流中使用动态配置
func DynamicConfigWorkflow(ctx workflow.Context) error {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 从配置源获取配置
    maxRetries := getConfigInt("max_retries", 3)

    // 使用配置
    retryPolicy := &temporal.RetryPolicy{
        MaximumAttempts: maxRetries,
    }

    ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute,
        RetryPolicy:         retryPolicy,
    })

    return workflow.ExecuteActivity(ctx, MyActivity, nil).Get(ctx, nil)
}

func getConfigInt(key string, defaultValue int) int {
    // 从动态配置源获取
    value := configSource.Get(key)
    if v, ok := value.(int); ok {
        return v
    }
    return defaultValue
}
```
