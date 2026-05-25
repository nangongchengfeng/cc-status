# 检测操作系统
ifeq ($(OS),Windows_NT)
    # Windows
    RM = powershell -Command Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
    CP = powershell -Command Copy-Item -Recurse -Force
    MKDIR = powershell -Command New-Item -ItemType Directory -Force | Out-Null
else
    # Linux/macOS
    RM = rm -rf
    CP = cp -r
    MKDIR = mkdir -p
endif

.PHONY: all build build-web build-server clean

# 默认目标
all: build

# 完整构建：前端 + 后端
build: build-web build-server

# 构建前端
build-web:
	@echo "Building web..."
	cd web && npm run build
	@echo "Copying web dist to server..."
	$(RM) server/internal/handler/ui/dist
	$(MKDIR) server/internal/handler/ui/dist
	$(CP) web/dist/* server/internal/handler/ui/dist/

# 构建后端（带 embed tag）
build-server:
	@echo "Building server..."
	cd server && go build -tags embed -o bin/server ./cmd/server

# 清理构建产物
clean:
	@echo "Cleaning..."
	$(RM) server/bin
	$(RM) server/internal/handler/ui/dist
	$(RM) web/dist
