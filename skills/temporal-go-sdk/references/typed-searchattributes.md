# Typed Search Attributes 类型化搜索属性

使用类型化的搜索属性。

```go
// 定义搜索属性
var (
    CustomerIDAttribute = searchattribute.SearchAttribute{
        Name:        "CustomerID",
        ValueType:   searchattribute.String,
        StringValue: pointer.String(""),
    }
    OrderCountAttribute = searchattribute.SearchAttribute{
        Name:        "OrderCount",
        ValueType:   searchattribute.Int,
        Int64Value:  pointer.Int64(0),
    }
    TotalAmountAttribute = searchattribute.SearchAttribute{
        Name:         "TotalAmount",
        ValueType:    searchattribute.Double,
        DoubleValue:  pointer.Float64(0.0),
    }
)

// 启动带搜索属性的工作流
func StartWithTypedSearchAttributes(ctx context.Context, client temporal.Client) error {
    searchAttributes := map[string]interface{}{
        "CustomerID":  "customer-123",
        "OrderCount": int64(5),
        "TotalAmount": 1500.50,
    }

    options := client.StartWorkflowOptions{
        ID:             "workflow-" + uuid.NewString(),
        TaskQueue:      "my-task-queue",
        SearchAttributes: searchAttributes,
    }

    _, err := client.StartWorkflow(ctx, options, MyWorkflow, input)
    return err
}

// 在工作流中获取搜索属性
func WorkflowWithSearchAttributes(ctx workflow.Context) error {
    info := workflow.GetInfo(ctx)

    // 获取搜索属性
    customerID := info.SearchAttributes["CustomerID"]
    orderCount := info.SearchAttributes["OrderCount"]

    return nil
}
```
