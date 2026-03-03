# Saga 分布式事务

使用 Saga 模式实现多步骤事务的补偿逻辑。

## Saga 结构

```go
type Saga struct {
    Activities []Activity
}

type Activity struct {
    Name string
    Compensate interface{}
}

func (s *Saga) Add(name string, compensate interface{}) {
    s.Activities = append(s.Activities, Activity{
        Name:       name,
        Compensate: compensate,
    })
}

func (s *Saga) Compensate(ctx workflow.Context) error {
    for i := len(s.Activities) - 1; i >= 0; i-- {
        activity := s.Activities[i]
        if activity.Compensate != nil {
            err := workflow.ExecuteActivity(ctx, activity.Compensate).Get(ctx, nil)
            if err != nil {
                // 记录补偿失败
                return err
            }
        }
    }
    return nil
}
```

## 使用 Saga

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    saga := &Saga{}

    // 执行操作并记录补偿
    err := workflow.ExecuteActivity(ctx, CreateOrder, order).Get(ctx, nil)
    if err != nil {
        return err
    }
    saga.Add("CreateOrder", CancelOrderActivity)

    err = workflow.ExecuteActivity(ctx, ReserveInventory, order.Items).Get(ctx, nil)
    if err != nil {
        saga.Compensate(ctx)
        return err
    }
    saga.Add("ReserveInventory", ReleaseInventoryActivity)

    err = workflow.ExecuteActivity(ctx, ProcessPayment, order.Payment).Get(ctx, nil)
    if err != nil {
        saga.Compensate(ctx)
        return err
    }
    saga.Add("ProcessPayment", RefundPaymentActivity)

    return nil
}
```

## 最佳实践

1. **反向执行**: 补偿按相反顺序执行
2. **幂等补偿**: 补偿逻辑必须幂等
3. **错误处理**: 记录补偿失败并告警
4. **部分成功**: 允许部分步骤成功
