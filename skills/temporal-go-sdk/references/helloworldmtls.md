# mTLS Hello World mTLS 问候

使用 mTLS 双向认证连接 Temporal Server。

```go
import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
)

// 加载 mTLS 证书
func LoadMTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
    // 加载客户端证书
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }

    // 加载 CA 证书
    caCert, err := ioutil.ReadFile(caFile)
    if err != nil {
        return nil, err
    }

    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:     caCertPool,
        MinVersion:   tls.VersionTLS12,
    }, nil
}

// 创建 mTLS Client
func CreateMTLSClient(ctx context.Context) (temporal.Client, error) {
    tlsConfig, err := LoadMTLSConfig("client.pem", "client-key.pem", "ca.pem")
    if err != nil {
        return nil, err
    }

    client, err := temporal.NewClient(temporal.ClientOptions{
        HostPort:  "localhost:7233",
        TLSConfig: tlsConfig,
    })

    return client, err
}
```
