# API Key Authentication API 密钥认证

使用 API Key 进行认证。

```go
// 创建 API Key 认证的 Client
func CreateAPIKeyClient(ctx context.Context, apiKey string) (temporal.Client, error) {
    client, err := temporal.NewClient(temporal.ClientOptions{
        HostPort: "localhost:7233",
        // API Key 认证
        Auth: &temporal.Auth{
            APIKey: apiKey,
        },
    })
    return client, err
}

// 在 Activity 中使用 API Key
func ActivityWithAPIKey(ctx context.Context, request Request) error {
    // 从 context 中获取 API Key
    apiKey := ctx.Value("api_key")

    client := &http.Client{}
    req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
    req.Header.Add("Authorization", "Bearer "+apiKey.(string))

    resp, err := client.Do(req)
    return err
}
```
