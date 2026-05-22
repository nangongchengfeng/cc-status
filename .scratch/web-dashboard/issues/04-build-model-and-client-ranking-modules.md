Status: done

# 落地模型排行与客户端排行模块

## What to build

在大屏页面中实现模型排行和客户端排行模块，并让它们消费 dashboard 接口中的 TOP 10 聚合结果。该切片需要把模型展示名回退逻辑、客户端总费用排行和横向图表或等效视觉呈现落地成可演示模块，使页面下半区具备明显的对比和分析价值。

## Acceptance criteria

- [x] 页面展示模型排行模块，固定渲染 dashboard 返回的 TOP 10 模型结果。
- [x] 模型排行优先显示 `display_name`，缺失时回退显示原始 `model`。
- [x] 页面展示客户端排行模块，固定渲染 dashboard 返回的 TOP 10 客户端费用结果。
- [x] 排行模块具备适合大屏的可视化呈现方式，能清晰比较不同项的高低。
- [x] 当排行数据为空时，模块仍展示稳定空态和标题结构。

## Blocked by

- `.scratch/web-dashboard/issues/02-connect-dashboard-query-and-time-range-presets.md`
