# Gin RESTful API 项目

小型 RESTful API 服务，Gin + GORM，分层架构。

## 项目架构

```
cmd/server/           # 入口：启动、依赖注入、路由注册
internal/
├── handler/          # HTTP 层：参数绑定、响应序列化
├── service/          # 业务逻辑层
├── repository/       # 数据访问层（GORM）
├── model/
│   ├── entity/       # GORM 模型（数据库映射）
│   └── dto/          # 请求/响应结构体
├── middleware/        # Gin 中间件
└── config/           # 配置加载
migrations/           # 数据库迁移
```

**分层约束：**
- handler 只做参数绑定和响应，不含业务逻辑
- service 之间可互相调用，repository 之间不可
- model 分为 entity（数据库映射）和 dto（请求/响应）
- `internal/` 下按层分包，不按业务域分包

## Gin 路由与中间件

**路由注册：**
- 集中在 `cmd/server/main.go` 或 `internal/handler/router.go`
- 按资源分组，前缀 `/api/v1`
- 每个资源一个 handler 文件，如 `user_handler.go`

**路由风格：**
- RESTful：`GET /users`、`POST /users`、`GET /users/:id`、`PUT /users/:id`、`DELETE /users/:id`
- URL 小写 + 连字符（`/user-profiles`，不用 `/userProfiles`）
- 路径参数 `:id`，查询参数 `?page=1&size=20`

**中间件顺序：** Recovery → Logger → CORS → 业务中间件
- 全局中间件：Logger、Recovery、CORS
- 业务中间件放 `internal/middleware/`，每个文件一个

**错误响应格式：**
```json
{ "code": "NOT_FOUND", "message": "用户不存在" }
```
- `code`：业务错误码，大写下划线
- `message`：中文描述
- HTTP 状态码与业务错误码分离

## GORM 约定

**模型定义：**
- 放 `internal/model/entity/`，每个资源一个文件
- 嵌入 `gorm.Model`
- 用 `TableName()` 显式指定表名
- 字段标签顺序：`gorm` → `json` → `validate`

**查询风格：**
- handler 中禁止直接写 GORM 查询，必须通过 repository
- repository 返回 entity 结构体，不返回 `*gorm.DB`
- 复杂查询用 scope 封装
- 分页封装：`Paginate(page, size) func(db *gorm.DB) *gorm.DB`

**事务：**
- 跨 repository 操作在 service 层用 `db.Transaction()` 包裹
- 禁止在 handler 层管理事务

**软删除：** 使用 GORM 默认软删除，硬删除用 `db.Unscoped().Delete()`

## API 设计

**请求：**
- 请求体用 dto 结构体 + `ShouldBindJSON` / `ShouldBindQuery`
- DTO 放 `internal/model/dto/`，文件按资源命名
- 必填字段 `binding:"required"`，路径参数用 `ShouldBindUri`

**响应：**
- 成功：`{ "data": { ... } }`
- 分页：`{ "data": [...], "total": 100, "page": 1, "size": 20 }`
- 响应 DTO 与 entity 分离，不暴露内部字段
- 用 `c.JSON(statusCode, response)` 返回，禁止在 handler 外设置响应

**状态码：**
- 200 成功 / 201 创建 / 204 删除（无响应体）
- 400 参数错误 / 401 未认证 / 403 无权限 / 404 不存在
- 409 冲突 / 500 内部错误
- 禁止用 200 + 业务错误码代替正确的 HTTP 状态码

## 测试

- `testing` + `testify`，测试文件与源文件同目录，`*_test.go`
- handler 测试：`httptest` 模拟请求
- service 测试：mock repository 接口
- repository 测试：SQLite 内存数据库集成测试
- 运行：`go test ./...`

## 配置

- `viper` 加载，优先级：环境变量 > 配置文件 > 默认值
- 配置结构体放 `internal/config/`
- 敏感信息只用环境变量

## 依赖注入

- 构造函数注入，不用 DI 框架
- service 定义自己需要的 repository 接口（接口定义在消费方）
