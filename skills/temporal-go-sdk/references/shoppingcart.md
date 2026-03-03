# Shopping Cart 购物车

购物车工作流示例。

```go
type CartItem struct {
    ProductID string
    Quantity  int
    Price     float64
}

type CartState struct {
    Items    []CartItem
    Total    float64
    Status   string
}

// 购物车工作流
func ShoppingCartWorkflow(ctx workflow.Context, userID string) (*CartState, error) {
    state := &CartState{
        Items:  []CartItem{},
        Total:  0,
        Status: "open",
    }

    // 注册添加商品信号
    workflow.SetSignalHandler(ctx, "addItem", func(item CartItem) {
        state.Items = append(state.Items, item)
        state.Total += item.Price * float64(item.Quantity)
    })

    // 注册移除商品信号
    workflow.SetSignalHandler(ctx, "removeItem", func(productID string) {
        for i, item := range state.Items {
            if item.ProductID == productID {
                state.Total -= item.Price * float64(item.Quantity)
                state.Items = append(state.Items[:i], state.Items[i+1:]...)
                break
            }
        }
    })

    // 注册结账信号
    workflow.SetSignalHandler(ctx, "checkout", func(_ string) {
        state.Status = "checking_out"
    })

    // 等待结账
    workflow.Await(ctx, func() bool { return state.Status == "checking_out" })

    // 处理结账
    err := workflow.ExecuteActivity(ctx, ProcessCheckoutActivity, state).Get(ctx, nil)
    if err != nil {
        state.Status = "checkout_failed"
        return state, err
    }

    state.Status = "completed"
    return state, nil
}
```
