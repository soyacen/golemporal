# Server JWT Auth 服务器 JWT 认证

使用 JWT Token 进行认证。

```go
import (
    "github.com/golang-jwt/jwt/v5"
)

// JWT 认证的 Client
func CreateJWTAuthClient(ctx context.Context, jwtToken string) (temporal.Client, error) {
    client, err := temporal.NewClient(temporal.ClientOptions{
        HostPort: "localhost:7233",
        Auth: &temporal.Auth{
            JWT: jwtToken,
        },
    })
    return client, err
}

// 生成 JWT Token
func GenerateJWTToken(claims jwt.Claims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("my-secret-key"))
}

// 在 Activity 中验证 JWT
func ActivityWithJWT(ctx context.Context, request Request) error {
    // 从 context 中获取 token
    token := ctx.Value("jwt_token")

    // 解析和验证 token
    parsedToken, err := jwt.Parse(token.(string), func(token *jwt.Token) (interface{}, error) {
        return []byte("my-secret-key"), nil
    })

    if err != nil || !parsedToken.Valid {
        return errors.New("invalid token")
    }

    // 处理请求
    return processRequest(request)
}
```
