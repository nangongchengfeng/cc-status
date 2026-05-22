Status: done

# 接入仪表盘查询与时间范围预设

## What to build

在前端中接入 dashboard 统计接口和最近请求接口的基础查询能力，并实现时间范围预设到显式 API 参数的转换。该切片需要打通 API 客户端、自定义 hooks、TanStack Query 和页面筛选状态，让前端能按“今天、最近 7 天、最近 30 天、本月、全部”这些预设请求服务端数据，同时保持展示格式化留在前端。

## Acceptance criteria

- [x] 前端已接入 dashboard 统计接口和最近请求接口的基础请求函数与查询 hooks。
- [x] 页面提供时间范围预设，并能把预设转换为 `start_at`、`end_at`、`interval` 请求参数。
- [x] 切换时间范围会触发相关查询重新获取数据，并保持 query key 与时间参数一致。
- [x] 前端不依赖后端返回展示文案，金额、时间和标签格式化逻辑仍保留在前端。
- [x] 当请求失败时，页面能获得稳定的错误状态供后续 UI 使用。

## Blocked by

- `.scratch/web-dashboard/issues/01-bootstrap-dashboard-web-app-and-home-route.md`
