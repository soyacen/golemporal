# Search Attributes 搜索属性

使用 Search Attributes 索引和检索工作流。

## 启动时添加 Search Attributes

```go
workflowOptions := client.StartWorkflowOptions{
    ID:        "my-workflow-id",
    TaskQueue: "my-task-queue",
    SearchAttributes: map[string]interface{}{
        "CustomStatus":    []string{"pending", "running"},
        "CustomerID":      []string{"customer-123"},
        "OrderCount":      []int{10},
        "Priority":         []int{1},
        "StartedAt":       []time.Time{time.Now()},
    },
}

client.ExecuteWorkflow(context.Background(), workflowOptions, MyWorkflow, input)
```

## 列表过滤

```go
// 过滤工作流
iter := c.ListWorkflow(context.Background(), &client.ListWorkflowsRequest{
    Query: "CustomStatus = 'pending' AND CustomerID = 'customer-123'",
})

for iter.Next() {
    fmt.Println(iter.Workflow().ID)
}
```

## Search Attributes 类型

| 类型 | 示例 |
|------|------|
| String | []string{"value"} |
| Int | []int{123} |
| Bool | []bool{true} |
| Datetime | []time.Time{time.Now()} |
| Double | []float64{1.5} |

## 在 Temporal UI 中使用

Search Attributes 会在 Temporal Web UI 中显示，可用于过滤工作流。

## 最佳实践

1. **预先定义**: 在 Temporal Server 中配置 Search Attributes
2. **索引优化**: 只索引常用查询字段
3. **复合查询**: 支持多字段组合查询
