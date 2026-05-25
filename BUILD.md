# 项目打包教程

本教程介绍如何打包 cc-status 项目的各个组件（client、server、web）。

## 目录

- [前置条件](#前置条件)
- [快速开始](#快速开始)
- [详细打包指南](#详细打包指南)
  - [打包 Client](#打包-client)
  - [打包 Server](#打包-server)
  - [打包 Web](#打包-web)
- [跨平台编译](#跨平台编译)
- [CI/CD 自动构建](#cicd-自动构建)

---

## 前置条件

在开始打包前，请确保已安装以下工具：

- **Go 1.25+** - 用于编译 client 和 server
- **Node.js 22+** - 用于编译 web 前端
- **Git** - 用于版本控制（可选）

### 验证安装

```powershell
# 验证 Go
go version

# 验证 Node.js
node --version
npm --version
```

---

## 快速开始

### 方法 1：使用 PowerShell 脚本（Windows 推荐）

```powershell
# 打包所有组件（web + server + client）
.\build.ps1 -Target all

# 只打包 client
.\build.ps1 -Target client

# 只打包 server
.\build.ps1 -Target server

# 只打包 web
.\build.ps1 -Target web

# 清理构建产物
.\build.ps1 -Target clean
```

### 方法 2：使用 Makefile（Linux/macOS/Windows（需安装 make）

```bash
# 打包所有组件（web + server）
make build

# 只打包 client
make build-client

# 只打包 server
make build-server

# 只打包 web
make build-web

# 清理构建产物
make clean
```

---

## 详细打包指南

### 打包 Client

Client 是一个命令行工具，用于采集和上报 Claude Code 使用数据。

#### 使用 PowerShell 脚本

```powershell
.\build.ps1 -Target client
```

输出文件：`client/bin/cc-usage-client.exe`（Windows）或 `client/bin/cc-usage-client`（Linux/macOS）

#### 使用 Makefile

```bash
make build-client
```

#### 手动使用 go build

```powershell
# 进入 client 目录
cd client

# Windows
$env:CGO_ENABLED=0
go build -ldflags="-s -w" -o bin/cc-usage-client.exe ./cmd/cc-usage-client

# Linux/macOS
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/cc-usage-client ./cmd/cc-usage-client
```

#### 构建参数说明

- `CGO_ENABLED=0` - 禁用 CGO，生成静态链接的可执行文件，不依赖系统库
- `-ldflags="-s -w"` - 去掉调试信息和符号表，减小文件体积（约减少 30-50%）
- `-o bin/xxx.exe` - 指定输出路径

#### 使用打包好的 Client

```powershell
# 查看帮助
.\client\bin\cc-usage-client.exe --help

# 试运行（不实际上报）
.\client\bin\cc-usage-client.exe dry-run

# 同步数据
.\client\bin\cc-usage-client.exe sync
```

注意：使用前需要先配置 `~/.cc-usage-client/config.yaml`，参考 `client/config.example.yaml`。

---

### 打包 Server

Server 是一个 Web 服务，提供 API 和数据统计功能。

#### 使用 PowerShell 脚本

```powershell
.\build.ps1 -Target server
```

输出文件：`server/bin/server.exe`（Windows）或 `server/bin/server`（Linux/macOS）

#### 使用 Makefile

```bash
make build-server
```

#### 手动使用 go build

```powershell
# 进入 server 目录
cd server

# Windows
$env:CGO_ENABLED=0
go build -tags embed -ldflags="-s -w" -o bin/server.exe ./cmd/server

# Linux/macOS
CGO_ENABLED=0 go build -tags embed -ldflags="-s -w" -extldflags=-static -o bin/server ./cmd/server
```

#### 构建参数说明

- `-tags embed` - 启用嵌入模式，将 web 前端静态资源打包进二进制文件
- `-extldflags=-static` - 静态链接（仅 Linux/macOS）

#### 使用打包好的 Server

```powershell
# 设置环境变量（必填）
$env:CC_USAGE_SERVER_AUTH_TOKEN="your-secret-token"

# 启动服务
.\server\bin\server.exe
```

访问 http://localhost:8080 查看仪表板。

可选环境变量：
- `CC_USAGE_SERVER_LISTEN_ADDR` - 监听地址，默认 `:8080`
- `CC_USAGE_SERVER_SQLITE_PATH` - SQLite 数据库路径，默认 `./server/data/server.db`

---

### 打包 Web

Web 是一个 React 前端应用，通常不需要单独打包，会在构建 Server 时自动嵌入。

#### 使用 PowerShell 脚本

```powershell
.\build.ps1 -Target web
```

输出目录：`web/dist/`

#### 使用 Makefile

```bash
make build-web
```

#### 手动使用 npm

```powershell
# 进入 web 目录
cd web

# 安装依赖（如果还没安装）
npm ci

# 构建
npm run build
```

构建产物会自动复制到 `server/internal/handler/ui/dist/`，供 Server 嵌入使用。

---

## 跨平台编译

Go 支持跨平台编译，你可以在 Windows 上编译 Linux/macOS 的可执行文件，反之亦然。

### 编译 Windows 版本（在任意平台）

```bash
# 编译 Windows amd64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o cc-usage-client-windows-amd64.exe ./client/cmd/cc-usage-client
```

### 编译 Linux 版本

```bash
# 编译 Linux amd64
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o cc-usage-client-linux-amd64 ./client/cmd/cc-usage-client
```

### 编译 macOS 版本

```bash
# 编译 macOS amd64 (Intel)
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o cc-usage-client-darwin-amd64 ./client/cmd/cc-usage-client

# 编译 macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o cc-usage-client-darwin-arm64 ./client/cmd/cc-usage-client
```

---

## CI/CD 自动构建

项目已经配置了 GitHub Actions 自动构建，当推送 tag 时会自动触发构建和发布。

### 触发自动构建

```bash
# 创建并推送 tag
git tag v1.0.0
git push origin v1.0.0
```

### CI/CD 流程

1. 运行测试（Windows/Linux）
2. 构建 Web 前端
3. 构建 Server（Windows/Linux，嵌入 Web）
4. 构建 Client（Windows/Linux）
5. 创建 GitHub Release，上传构建产物

### 查看构建结果

访问项目的 GitHub Actions 页面：`https://github.com/nangongchengfeng/cc-status/actions`

---

## 目录结构

构建完成后，项目结构如下：

```
cc-status/
├── client/
│   └── bin/
│       └── cc-usage-client.exe    # Client 可执行文件
├── server/
│   └── bin/
│       └── server.exe               # Server 可执行文件（包含 Web 前端）
└── web/
    └── dist/                       # Web 构建产物（已嵌入 Server）
```

---

## 常见问题

### Q: 构建的 exe 文件太大怎么办？

A: 使用 `-ldflags="-s -w"` 参数可以显著减小文件体积。如果还需要进一步减小，可以使用 UPX 压缩工具：

```powershell
upx --best --lzma .\client\bin\cc-usage-client.exe
```

### Q: 如何在老版本 Linux（如 CentOS 7）上运行？

A: 使用 `CGO_ENABLED=0` 和 `-extldflags=-static` 参数构建静态链接的可执行文件，不依赖系统 GLIBC 版本。

### Q: 构建 Server 时提示找不到 `ui/dist`？

A: 需要先构建 Web 前端，或者确保 `server/internal/handler/ui/dist/` 目录存在（至少有 `.gitkeep` 文件）。
