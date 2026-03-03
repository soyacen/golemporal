# Workflow Security Interceptor 工作流安全拦截器

拦截工作流调用以添加安全检查。

```go
// 安全拦截器
type SecurityInterceptor struct {
    allowedWorkflows map[string]bool
    allowedUsers    map[string][]string
}

func NewSecurityInterceptor() *SecurityInterceptor {
    return &SecurityInterceptor{
        allowedWorkflows: map[string]bool{
            "MyWorkflow":      true,
            "AdminWorkflow":   true,
        },
        allowedUsers: map[string][]string{
            "admin": {"MyWorkflow", "AdminWorkflow"},
            "user":  {"MyWorkflow"},
        },
    }
}

func (s *SecurityInterceptor) InterceptWorkflow(ctx workflow.Context,
    method string, args []interface{}, next workflow.NextInterceptor) (interface{}, error) {

    // 获取调用者信息
    info := workflow.GetInfo(ctx)
    caller := info.ExternalOwner

    // 检查权限
    if !s.checkPermission(caller, method) {
        return nil, fmt.Errorf("permission denied for workflow: %s", method)
    }

    // 继续执行
    return next.Execute(ctx, args...)
}

func (s *SecurityInterceptor) checkPermission(caller, workflowName string) bool {
    workflows, ok := s.allowedUsers[caller]
    if !ok {
        return false
    }

    for _, wf := range workflows {
        if wf == workflowName {
            return true
        }
    }

    return false
}

// 注册拦截器
func init() {
    workflow.RegisterInterceptable(new(SecurityInterceptor))
}
```
