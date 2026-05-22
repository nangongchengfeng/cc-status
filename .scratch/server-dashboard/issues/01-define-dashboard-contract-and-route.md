Status: done

# 定义仪表盘统计契约并接入路由

## What to build

新增一个面向 Web 数据大屏的独立统计接口，并把它接入现有服务端路由体系。该切片需要打通从 HTTP 路由、查询参数绑定、响应包装到基础 service 契约的完整链路，但不要求在本 issue 中完成所有聚合细节。目标是让前端和后续切片有一个稳定的 dashboard API 入口，并明确它与现有 `overview`、`trend`、`logs` 的职责边界。

## Acceptance criteria

- [x] 服务端新增受静态 Bearer token 保护的 dashboard 统计接口，并接入现有 `/api/v1` 路由体系。
- [x] 接口支持显式 `start_at`、`end_at`、`interval` 查询参数，且 `interval` 仅允许 `hour` 或 `day`。
- [x] 成功响应使用查询类接口统一的 `{ "data": ... }` 包装，非法参数返回现有风格的错误响应。
- [x] 接口返回结构已预留总览卡片、统一时间桶趋势、模型排行、客户端排行、缓存效益分析这些大屏数据域。
- [x] 新接口不替代现有 `overview`、`trend`、`logs`，最近请求列表仍明确由 `/api/v1/logs` 提供。

## Blocked by

None - can start immediately
