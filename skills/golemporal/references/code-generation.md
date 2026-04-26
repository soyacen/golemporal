# 代码生成参考

本参考文档说明如何安装 golemporal 插件、运行代码生成命令，以及理解生成产物。

## 前置依赖

### 1. 安装 protoc

从 [protobuf releases](https://github.com/protocolbuffers/protobuf/releases) 下载对应系统的 `protoc` 编译器，或：

```bash
# macOS
brew install protobuf

# Ubuntu/Debian
apt-get install -y protobuf-compiler

# 验证安装
protoc --version
```

### 2. 安装 Go protobuf 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

确保 `$GOPATH/bin`（通常是 `~/go/bin`）在 `PATH` 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 3. 安装 golemporal 插件

```bash
# 从源码安装（在 golemporal 仓库根目录）
go install ./cmd/protoc-gen-golemporal

# 或通过 go install 安装特定版本
go install github.com/soyacen/golemporal/cmd/protoc-gen-golemporal@v0.3.0
```

验证安装：

```bash
protoc-gen-golemporal --version
# 输出: protoc-gen-golemporal v0.3.0
```

## protoc 命令详解

### 基本命令

```bash
protoc \
  --proto_path=. \
  --go_out=. \
  --go_opt=paths=source_relative \
  --golemporal_out=. \
  --golemporal_opt=paths=source_relative \
  example.proto
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `--proto_path=.` | proto 文件的搜索路径，可多次指定 |
| `--go_out=.` | go protobuf 插件输出目录 |
| `--go_opt=paths=source_relative` | 输出文件与 proto 文件保持相对路径 |
| `--golemporal_out=.` | golemporal 插件输出目录 |
| `--golemporal_opt=paths=source_relative` | 输出文件与 proto 文件保持相对路径 |

### 多路径导入

当 proto 文件引用了其他目录中的 proto（如 golemporal 的 `metadata.proto`）时，需要添加额外的 `--proto_path`：

```bash
protoc \
  --proto_path=. \
  --proto_path=../../ \
  --go_out=. \
  --go_opt=paths=source_relative \
  --golemporal_out=. \
  --golemporal_opt=paths=source_relative \
  api/*.proto
```

`--proto_path=../../` 使 protoc 能够找到项目根目录下的 protobuf 定义。

### 使用脚本

项目中通常提供一个 `protoc.sh` 脚本：

```bash
#!/bin/bash

protoc \
--proto_path=. \
--proto_path=../../ \
--go_out=. \
--go_opt=paths=source_relative \
--golemporal_out=. \
--golemporal_opt=paths=source_relative \
*/*.proto
```

使用方式：

```bash
cd example && ./protoc.sh
```

## 生成产物详解

假设 proto 文件 `hello.proto` 定义了 `HelloWorkflow` 和 `GreetActivity`：

### 1. 标准 protobuf 生成

`hello.pb.go`：由 `protoc-gen-go` 生成，包含 message 的 Go struct 定义。

### 2. golemporal 生成

`hello_temporal.pb.go`：由 `protoc-gen-golemporal` 生成，包含以下内容：

#### Activity 相关

```go
// 类型常量
var GreetActivity_Greet_ActitityType = "/myapp.api.GreetActivity/Greet"

// Activity Client 接口
type GreetActivityClient interface {
    Greet(ctx workflow.Context, in *HelloRequest) (*HelloResponse, error)
}

// Activity Client 构造函数
func NewGreetActivityClient() GreetActivityClient

// Activity Server 接口
type GreetActivityServer interface {
    Greet(context.Context, *HelloRequest) (*HelloResponse, error)
}
```

#### Workflow 相关

```go
// 类型常量
var HelloWorkflow_Hello_WorkflowType = "/myapp.api.HelloWorkflow/Hello"

// Workflow Client 接口
type HelloWorkflowClient interface {
    Hello(ctx context.Context, in *HelloRequest, opts ...starter.Option) (*HelloResponse, *protobuf.Metadata, error)
}

// Workflow Client 构造函数
func NewHelloWorkflowClient(c client.Client, taskQueue string) HelloWorkflowClient

// Workflow Server 接口
type HelloWorkflowServer interface {
    Hello(workflow.Context, *HelloRequest) (*HelloResponse, error)
}
```

#### 注册函数

```go
func RegisterGreetActivity(wk worker.Worker, server GreetActivityServer) {
    wk.RegisterActivityWithOptions(server.Greet, activity.RegisterOptions{
        Name:                          GreetActivity_Greet_ActitityType,
        DisableAlreadyRegisteredCheck: true,
    })
}

func RegisterHelloWorkflow(wk worker.Worker, server HelloWorkflowServer) {
    wk.RegisterWorkflowWithOptions(server.Hello, workflow.RegisterOptions{
        Name:                          HelloWorkflow_Hello_WorkflowType,
        DisableAlreadyRegisteredCheck: true,
    })
}
```

## 重新生成代码的时机

需要在以下情况重新运行代码生成：

| 场景 | 是否需要重新生成 |
|------|----------------|
| 新增/删除/修改 proto message | 是 |
| 新增/删除/修改 proto service | 是 |
| 新增/删除/修改 RPC 方法 | 是 |
| 修改 `go_package` | 是 |
| 修改 workflow/activity 实现（.go 文件） | 否 |
| 更新 golemporal 插件版本 | 是（建议） |
| 更新 protoc 版本 | 视情况 |

## 常见问题

### Q: `protoc-gen-golemporal: program not found`

确保 `protoc-gen-golemporal` 在 `PATH` 中：

```bash
which protoc-gen-golemporal
# 如果未找到，添加 GOPATH/bin 到 PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Q: 生成的代码没有包含某些 service

检查 service 名称是否以 `Workflow` 或 `Activity` 结尾（区分大小写）。

### Q: 导入路径错误

确保 `go_package` 选项正确设置：

```protobuf
option go_package = "github.com/myorg/myapp/api;api";
```
