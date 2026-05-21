Status: done

# 搭建 server 运行时骨架与健康检查

## Type

AFK

## What to build

交付一个可启动的 Go server 基础骨架，能够读取环境变量配置、校验必填静态 token、初始化默认 SQLite 路径并启动 HTTP 服务。该切片需要打通全局中间件、统一错误响应、Bearer 鉴权中间件和 `GET /healthz` 存活探针。完成后，服务端可以在未接入业务 API 的情况下启动、拒绝未授权请求，并暴露无需鉴权的健康检查接口。

## Acceptance criteria

- [x] 服务端存在可运行入口，支持监听地址、SQLite 路径和静态 token 的环境变量配置
- [x] 静态 token 为必填配置，缺失时服务端启动失败并输出明确错误
- [x] `GET /healthz` 无需鉴权即可返回成功响应，其他受保护路由在缺少或错误 token 时返回 `401`
- [x] 统一注册基础中间件与错误响应约定，后续业务接口可直接复用

## Blocked by

None - can start immediately
