Status: done

# 初始化 client 运行时与本地状态

## What to build

交付一个可运行的 Go client 基础骨架，能够在 `client` 目录下启动单次 CLI，并完成 `~/.cc-usage-client` 的基础初始化。该切片需要打通配置加载、环境变量覆盖、本地 SQLite 初始化、稳定**客户端标识**生成与持久化，以及进程级互斥锁。完成后，`sync` 和 `dry-run` 两个命令都可以完成启动、初始化和退出，即使业务同步逻辑还未完全接入。

## Acceptance criteria

- [x] `client` 目录下存在可运行的 Go CLI 骨架，并提供 `sync` 与 `dry-run` 两个命令入口
- [x] 首次运行会创建 `~/.cc-usage-client`，并初始化 SQLite、`sync_state`、`reported_ids` 及保存**客户端标识**所需的元数据
- [x] `config.yaml` 为主配置来源，环境变量可以覆盖服务端 URL、认证 token、批大小和超时时间
- [x] 同一台机器上并发启动第二个实例时，会因进程级互斥锁而失败退出，并输出可理解的错误信息

## Blocked by

None - can start immediately
