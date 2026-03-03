# Encryption 数据加密

使用 DataConverter 加密工作流和 Activity 的数据。

## 安装加密包

```bash
go get go.temporal.io/sdk/contrib/data-converter/encryption
```

## 配置加密 DataConverter

```go
import "go.temporal.io/sdk/contrib/data-converter/encryption"

encryptionKey, err := base64.StdEncoding.DecodeString("my-32-byte-base64-encoded-key!!")
if err != nil {
    log.Fatal(err)
}

dataConverter := encryption.New(encryptionKey)

client, err := client.Dial(client.Options{
    DataConverter: dataConverter,
})

worker := worker.New(c, "task-queue", worker.Options{
    DataConverter: dataConverter,
})
```

## 使用自定义密钥 ID

```go
dataConverter := encryption.NewWithKeys(encryption.Key{
    ID:   "key-v1",
    Key:  encryptionKey,
})
```

## 最佳实践

1. **密钥管理**: 使用安全的密钥管理系统（如 Vault）
2. **密钥轮换**: 支持多个密钥 ID 实现无缝轮换
3. **性能**: 加密有性能开销，考虑只加密敏感字段
