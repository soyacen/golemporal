# Safe Message Handler 安全消息处理

安全处理外部消息。

```go
// 安全的消息处理器
type SafeMessageHandler struct {
    maxMessageSize int
    allowedTypes  map[string]bool
}

func NewSafeMessageHandler() *SafeMessageHandler {
    return &SafeMessageHandler{
        maxMessageSize: 1024 * 1024, // 1MB
        allowedTypes: map[string]bool{
            "user.created":   true,
            "user.updated":   true,
            "order.completed": true,
        },
    }
}

// 验证消息
func (h *SafeMessageHandler) ValidateMessage(msg Message) error {
    // 检查大小
    if len(msg.Data) > h.maxMessageSize {
        return errors.New("message too large")
    }

    // 检查类型
    if !h.allowedTypes[msg.Type] {
        return errors.New("invalid message type")
    }

    // 检查签名
    if !h.verifySignature(msg) {
        return errors.New("invalid signature")
    }

    return nil
}

// 处理消息
func (h *SafeMessageHandler) HandleMessage(ctx context.Context, msg Message) error {
    if err := h.ValidateMessage(msg); err != nil {
        return err
    }

    // 处理消息
    return processMessage(msg)
}
```
