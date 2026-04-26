# Proto 定义规范

本参考文档说明如何为 golemporal 定义 proto 文件，包括 service 命名约定、message 定义和 go_package 配置。

## 完整示例

```protobuf
syntax = "proto3";

package myapp.order.api;

option go_package = "github.com/myorg/myapp/order/api;api";

// ========== Messages ==========

message CreateOrderRequest {
  string user_id = 1;
  repeated OrderItem items = 2;
  string currency = 3;
}

message OrderItem {
  string product_id = 1;
  int32 quantity = 2;
  double price = 3;
}

message CreateOrderResponse {
  string order_id = 1;
  double total_amount = 2;
  string status = 3;
}

message ChargeRequest {
  string order_id = 1;
  double amount = 2;
  string currency = 3;
}

message ChargeResponse {
  string transaction_id = 1;
  string status = 2;
}

message ShipRequest {
  string order_id = 1;
  string address = 2;
}

message ShipResponse {
  string tracking_id = 1;
}

// ========== Workflow Services ==========

service OrderWorkflow {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
}

// ========== Activity Services ==========

service PaymentActivity {
  rpc Charge(ChargeRequest) returns (ChargeResponse);
}

service ShippingActivity {
  rpc Ship(ShipRequest) returns (ShipResponse);
}
```

## Service 命名规则

golemporal 代码生成器通过 service 名称后缀识别服务类型（区分大小写）：

| 后缀 | 生成类型 | 示例 |
|------|---------|------|
| `Workflow` | Workflow 客户端、服务器、注册函数 | `OrderWorkflow` → `OrderWorkflowClient`、`OrderWorkflowServer`、`RegisterOrderWorkflow` |
| `Activity` | Activity 客户端、服务器、注册函数 | `PaymentActivity` → `PaymentActivityClient`、`PaymentActivityServer`、`RegisterPaymentActivity` |

### 重要规则

1. **精确后缀匹配**：必须是 `Workflow` 或 `Activity` 结尾，不能是 `workflow` 或 `activity`（小写不匹配）
2. **支持多服务**：一个 proto 文件可以包含多个 `Workflow` 和 `Activity` 服务
3. **无后缀服务**：不以 `Workflow` 或 `Activity` 结尾的 service 会被忽略，不生成任何代码

### 命名建议

```protobuf
// 推荐：使用业务领域 + 类型后缀
service OrderWorkflow { ... }      // 处理订单的工作流
service PaymentActivity { ... }    // 支付相关的 activity
service InventoryActivity { ... }  // 库存相关的 activity

// 不推荐：模糊或过于通用的命名
service MyWorkflow { ... }         // 语义不清
service DoActivity { ... }         // 语义不清
```

## Message 定义最佳实践

### 1. 请求/响应对应 RPC

每个 RPC 方法应有一对独立的 request/response message，即使字段相同：

```protobuf
// 推荐：每个 RPC 有独立的 message
service UserWorkflow {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
}

// 不推荐：多个 RPC 共享同一个 message
service UserWorkflow {
  rpc CreateUser(UserRequest) returns (UserResponse);  // 模糊
  rpc UpdateUser(UserRequest) returns (UserResponse);  // 难以区分
}
```

### 2. 使用 proto3 字段类型

```protobuf
message Example {
  string name = 1;           // 字符串
  int32 count = 2;           // 整数
  int64 timestamp = 3;       // 长整数（时间戳）
  double amount = 4;         // 浮点数
  bool active = 5;           // 布尔值
  repeated string tags = 6;  // 数组
  bytes data = 7;            // 二进制数据
}
```

### 3. 嵌套 Message

```protobuf
message OrderRequest {
  string user_id = 1;

  message Item {
    string product_id = 1;
    int32 quantity = 2;
  }

  repeated Item items = 2;
}
```

## go_package 选项

`go_package` 选项决定生成代码的导入路径：

```protobuf
option go_package = "github.com/myorg/myapp/api;api";
```

格式：`"import/path;package_name"`

- **import path**：Go 模块中该包的导入路径
- **package name**：生成代码的包名（如 `api`、`pb`、`workflowapi`）

### 示例配置

| 项目结构 | go_package |
|----------|-----------|
| `api/order.proto` | `"github.com/myorg/myapp/api;api"` |
| `proto/workflow.proto` | `"github.com/myorg/myapp/proto;proto"` |
| `internal/pb/service.proto` | `"github.com/myorg/myapp/internal/pb;pb"` |

## 多文件组织

对于大型项目，可以按业务领域拆分到多个 proto 文件：

```
api/
├── order.proto      # OrderWorkflow + PaymentActivity + ShippingActivity
├── user.proto       # UserWorkflow + EmailActivity
└── inventory.proto  # InventoryWorkflow + StockActivity
```

每个文件独立生成：

```bash
protoc --go_out=. --golemporal_out=. api/*.proto
```

生成：
- `order_temporal.pb.go`
- `user_temporal.pb.go`
- `inventory_temporal.pb.go`
