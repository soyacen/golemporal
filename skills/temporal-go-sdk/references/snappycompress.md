# Snappy Compression Snappy 压缩

使用 Snappy 压缩数据。

```go
import (
    "github.com/golang/snappy"
)

// 压缩 Activity 输入/输出
func CompressActivity(ctx context.Context, data []byte) ([]byte, error) {
    compressed := snappy.Encode(nil, data)
    return compressed, nil
}

func DecompressActivity(ctx context.Context, compressed []byte) ([]byte, error) {
    decompressed, err := snappy.Decode(nil, compressed)
    return decompressed, err
}

// 在工作流中使用压缩
func CompressedWorkflow(ctx workflow.Context, largeData []byte) error {
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 压缩数据
    var compressed []byte
    err := workflow.ExecuteActivity(ctx, CompressActivity, largeData).Get(ctx, &compressed)
    if err != nil {
        return err
    }

    // 传输压缩数据
    var result []byte
    err = workflow.ExecuteActivity(ctx, ProcessCompressedActivity, compressed).Get(ctx, &result)
    if err != nil {
        return err
    }

    // 解压结果
    var decompressed []byte
    return workflow.ExecuteActivity(ctx, DecompressActivity, result).Get(ctx, &decompressed)
}
```
