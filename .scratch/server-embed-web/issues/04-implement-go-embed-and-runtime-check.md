Status: done

# 实现 go:embed 与运行时检测

## What to build

使用 build tag 方案实现：ui.go（开发版，不 embed）和 ui_embed.go（构建版，使用 go:embed），开发阶段使用本地文件系统，构建阶段用 -tags embed 启用嵌入模式。

## Acceptance criteria

- [x] 使用 build tag 区分开发/embed 模式。
- [x] 运行时检测逻辑正确：无文件时不注册 UI 路由。
- [x] 开发阶段（无 dist 文件）server 仍能正常启动，API 正常工作。

## Blocked by

- 01-setup-ui-handler-framework
- 03-add-gitignore-and-dir-prep
