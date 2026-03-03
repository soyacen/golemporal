# PSO (Periodic Sampling Optimization) 周期采样优化

周期性采样执行任务。

```go
func PSOWorkflow(ctx workflow.Context, config PSOConfig) error {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 周期性采样
    ticker := workflow.NewTicker(ctx, config.SampleInterval)

    var samplesCollected int
    for {
        select {
        case <-ticker.Ch():
            // 采样执行
            if samplesCollected >= config.MaxSamples {
                ticker.Stop()
                return nil
            }

            var sampleResult SampleResult
            err := workflow.ExecuteActivity(ctx, SampleActivity, config).Get(ctx, &sampleResult)
            if err != nil {
                return err
            }

            samplesCollected++

            // 检查终止条件
            if sampleResult.Terminate {
                ticker.Stop()
                return nil
            }

        case <-ctx.Done():
            ticker.Stop()
            return ctx.Err()
        }
    }
}
```
