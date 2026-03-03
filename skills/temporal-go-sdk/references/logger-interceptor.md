# Logger Interceptor 日志拦截器

拦截工作流和 Activity 调用以添加日志。

```go
// 日志拦截器
type LoggerInterceptor struct {
    logger log.Logger
}

func (l *LoggerInterceptor) InterceptActivity(ctx context.Context, method string,
    req interface{}, info *activity.Info, next activity.NextInterceptor) (interface{}, error) {

    l.logger.Info("Activity started", "method", method)
    start := time.Now()

    result, err := next.Execute(ctx, req)

    l.logger.Info("Activity completed",
        "method", method,
        "duration", time.Since(start),
        "error", err)

    return result, err
}

func (l *LoggerInterceptor) InterceptWorkflow(ctx workflow.Context,
    method string, args []interface{}, next workflow.NextInterceptor) (interface{}, error) {

    l.logger.Info("Workflow started", "method", method)
    start := time.Now()

    result, err := next.Execute(ctx, args...)

    l.logger.Info("Workflow completed",
        "method", method,
        "duration", time.Since(start),
        "error", err)

    return result, err
}
```
