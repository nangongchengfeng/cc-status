# Client Bootstrap Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 在 `client` 目录下交付首个可运行的 Go CLI 骨架，完成应用数据目录、配置加载、SQLite 初始化、稳定客户端标识和进程级互斥锁。

**Architecture:** 采用小而深的运行时模块设计。CLI 只负责参数分发，运行时模块负责应用目录、配置和依赖初始化，状态存储模块负责 SQLite 与客户端标识，锁模块负责单实例运行。测试优先覆盖外部行为：初始化结果、配置覆盖和并发启动失败。

**Tech Stack:** Go、标准库、`modernc.org/sqlite`、`gopkg.in/yaml.v3`

---

### Task 1: 搭建最小 CLI 骨架

**Files:**
- Create: `client/go.mod`
- Create: `client/cmd/cc-usage-client/main.go`
- Create: `client/internal/cli/app.go`
- Test: `client/internal/cli/app_test.go`

**Step 1: Write the failing test**

编写一个通过公共入口运行 CLI 的测试，验证 `sync` 和 `dry-run` 子命令存在，并且在最小依赖注入下可以返回成功。

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run TestAppRunsKnownCommands -v`
Expected: FAIL，因为 CLI 入口尚不存在。

**Step 3: Write minimal implementation**

实现最小 `App` 结构和命令分发，让 `sync` 与 `dry-run` 能执行占位运行逻辑。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run TestAppRunsKnownCommands -v`
Expected: PASS

### Task 2: 初始化应用数据目录与 SQLite

**Files:**
- Create: `client/internal/runtime/runtime.go`
- Create: `client/internal/storage/sqlite.go`
- Test: `client/internal/runtime/runtime_test.go`

**Step 1: Write the failing test**

编写运行时初始化测试，验证首次运行会创建 `~/.cc-usage-client` 等价的测试目录、SQLite 文件和基础表结构。

**Step 2: Run test to verify it fails**

Run: `go test ./internal/runtime -run TestBootstrapCreatesAppDataDirAndDatabase -v`
Expected: FAIL，因为初始化逻辑尚不存在。

**Step 3: Write minimal implementation**

实现应用目录初始化、SQLite 打开与 `sync_state`、`reported_ids`、`metadata` 表创建。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/runtime -run TestBootstrapCreatesAppDataDirAndDatabase -v`
Expected: PASS

### Task 3: 生成并持久化客户端标识

**Files:**
- Modify: `client/internal/runtime/runtime.go`
- Modify: `client/internal/storage/sqlite.go`
- Test: `client/internal/runtime/runtime_test.go`

**Step 1: Write the failing test**

新增测试，验证首次初始化生成 UUID 形式的客户端标识，再次初始化返回同一个值。

**Step 2: Run test to verify it fails**

Run: `go test ./internal/runtime -run TestBootstrapPersistsStableClientID -v`
Expected: FAIL，因为客户端标识尚未持久化。

**Step 3: Write minimal implementation**

在 `metadata` 中保存 `client_id`，首次缺失时生成 UUID，后续直接复用。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/runtime -run TestBootstrapPersistsStableClientID -v`
Expected: PASS

### Task 4: 加入配置文件和环境变量覆盖

**Files:**
- Create: `client/internal/config/config.go`
- Modify: `client/internal/runtime/runtime.go`
- Test: `client/internal/config/config_test.go`

**Step 1: Write the failing test**

编写配置加载测试，验证 `config.yaml` 被读取，且环境变量可以覆盖 `server_url`、`auth_token`、`batch_size`、`timeout_seconds`。

**Step 2: Run test to verify it fails**

Run: `go test ./internal/config -run TestLoadConfigWithEnvOverride -v`
Expected: FAIL，因为配置模块尚不存在。

**Step 3: Write minimal implementation**

实现 YAML 配置解析、默认值填充和环境变量覆盖。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/config -run TestLoadConfigWithEnvOverride -v`
Expected: PASS

### Task 5: 加入单实例互斥锁并接入 CLI

**Files:**
- Create: `client/internal/lock/file_lock.go`
- Modify: `client/internal/runtime/runtime.go`
- Modify: `client/internal/cli/app.go`
- Test: `client/internal/lock/file_lock_test.go`
- Test: `client/internal/cli/app_test.go`

**Step 1: Write the failing test**

编写锁测试与 CLI 测试，验证第二个实例无法获取锁，并返回可理解错误。

**Step 2: Run test to verify it fails**

Run: `go test ./internal/lock ./internal/cli -run Test -v`
Expected: FAIL，因为锁逻辑尚未实现。

**Step 3: Write minimal implementation**

实现基于锁文件的单实例互斥，并在 `sync`/`dry-run` 启动时统一接入。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/lock ./internal/cli -run Test -v`
Expected: PASS

### Task 6: 验证、更新 issue 并提交

**Files:**
- Modify: `.scratch/claude-usage-sync/issues/01-bootstrap-client-runtime-and-state.md`

**Step 1: Run focused verification**

Run: `go test ./...`
Expected: PASS

**Step 2: Update issue**

将本 issue 的复选框全部改为已完成，并把 `Status:` 更新为已完成状态。

**Step 3: Commit**

Run: `git add client docs/plans .scratch/claude-usage-sync/issues/01-bootstrap-client-runtime-and-state.md`

Commit message: `feat(client): 初始化运行时与本地状态`
