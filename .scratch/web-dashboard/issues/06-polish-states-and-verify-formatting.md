Status: done

# 补齐空态错误态与展示格式化验证

## What to build

对大屏前端做最后一轮状态打磨和关键格式化验证，确保页面在 loading、error、empty、success 四类状态下都有稳定表现。该切片需要补齐金额、时间、长标签等前端展示格式化的收口逻辑，并通过最小自动化测试或构建验证固定这些对用户可观察的行为。

## Acceptance criteria

- [x] 页面在 loading、error、empty、success 状态下都有可观察且稳定的展示。
- [x] 金额、时间和长标签的前端格式化逻辑已经统一，且不会污染后端 API 契约。
- [x] 至少为关键数据转换逻辑或关键页面状态补齐自动化测试，或提供明确的最小构建验证证据。
- [x] 大屏页面在主要模块都接入完成后仍保持布局稳定，不因边界数据出现明显错位。
- [x] 该切片完成后，前端具备可演示的整体完成度。

## Blocked by

- `.scratch/web-dashboard/issues/03-build-overview-cards-and-trend-section.md`
- `.scratch/web-dashboard/issues/04-build-model-and-client-ranking-modules.md`
- `.scratch/web-dashboard/issues/05-build-cache-benefit-and-recent-requests.md`
