# mTLS 认证

使用 mTLS 连接到 Temporal Cloud。

## 基本 mTLS 配置

```go
c, err := client.Dial(client.Options{
    HostPort:  "my-cluster.tmprl.cloud:7233",
    Namespace: "my-namespace",
    ConnectionOptions: client.ConnectionOptions{
        TLS: &tls.Config{
            CertFile: "cert.pem",
            KeyFile:  "key.pem",
            RootCAFile: "ca.pem",
        },
    },
})
```

## 动态 mTLS 凭证刷新

参考: dynamicmtls 示例

```go
// 实现凭证刷新器
type CredentialRefresher struct {
    certPath string
    keyPath  string
}

func (r *CredentialRefresher) GetCurrentCertificate() (*tls.Certificate, error) {
    cert, err := tls.LoadX509KeyPair(r.certPath, r.keyPath)
    return &cert, err
}

// 配置动态 TLS
tlsConfig := &tls.Config{
    GetClientCertificate: credentialRefresher.GetCurrentCertificate,
}
```

## Worker 配置

```go
w := worker.New(c, "task-queue", worker.Options{
    ConnectionOptions: worker.ConnectionOptions{
        TLS: &tls.Config{
            CertFile: "cert.pem",
            KeyFile:  "key.pem",
            RootCAFile: "ca.pem",
        },
    },
})
```
