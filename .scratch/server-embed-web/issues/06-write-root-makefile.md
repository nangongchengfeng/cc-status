Status: done

# 编写根目录 Makefile

## What to build

在仓库根目录创建 Makefile，包含 build-web、build-server、build 三个主要目标。build-web 负责构建前端并跨平台复制到 server/internal/ui/dist。

## Acceptance criteria

- [x] 根目录 Makefile 已创建。
- [x] make build-web 能正确执行 npm run build 并复制到 server/internal/ui/dist。
- [x] make build-server 能正确构建 server 二进制。
- [x] make build 能依次执行 build-web 和 build-server。
- [x] 复制逻辑在 Windows 和 Linux/macOS 都能工作（检测 OS 使用不同命令）。

## Blocked by

- 03-add-gitignore-and-dir-prep
- 04-implement-go-embed-and-runtime-check
