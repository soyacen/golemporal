# Batch Sliding Window 滑动窗口批处理

滑动窗口模式的批处理。

```go
type BatchWindow struct {
    items     []Item
    size      int
    mu        sync.Mutex
}

func NewBatchWindow(size int) *BatchWindow {
    return &BatchWindow{
        items: make([]Item, 0, size),
        size:  size,
    }
}

// 添加项目
func (b *BatchWindow) Add(item Item) bool {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.items = append(b.items, item)
    return len(b.items) >= b.size
}

// 获取批次
func (b *BatchWindow) GetBatch() []Item {
    b.mu.Lock()
    defer b.mu.Unlock()

    batch := b.items
    b.items = make([]Item, 0, b.size)
    return batch
}

// 滑动窗口工作流
func SlidingWindowWorkflow(ctx workflow.Context, items []Item) error {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    window := NewBatchWindow(10)

    // 处理每个项目
    for _, item := range items {
        isFull := window.Add(item)

        if isFull || len(items) == 0 {
            // 获取当前批次
            batch := window.GetBatch()

            // 处理批次
            err := workflow.ExecuteActivity(ctx, ProcessBatchActivity, batch).Get(ctx, nil)
            if err != nil {
                return err
            }
        }
    }

    // 处理剩余项目
    remaining := window.GetBatch()
    if len(remaining) > 0 {
        return workflow.ExecuteActivity(ctx, ProcessBatchActivity, remaining).Get(ctx, nil)
    }

    return nil
}
```
