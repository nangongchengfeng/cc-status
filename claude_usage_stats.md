# Claude Code 使用统计系统详解

## 目录

1. [系统概述](#系统概述)
2. [核心架构](#核心架构)
3. [Claude 会话日志解析](#claude-会话日志解析)
4. [Codex 会话日志解析](#codex-会话日志解析)
5. [Gemini 会话日志解析](#gemini-会话日志解析)
6. [去重机制](#去重机制)
7. [同步状态跟踪](#同步状态跟踪)
8. [数据存储](#数据存储)
9. [数据库 Schema 参考](#数据库-schema-参考)

---

## 系统概述

### 功能简介

从 Claude Code、Codex、Gemini CLI 的本地会话日志文件中提取 token 使用数据，实现无需代理拦截的使用统计功能（v3.13 新增）。

### 数据流

```
本地会话文件 → 增量解析/去重 → 费用计算 → proxy_request_logs 表
                    ↑
              session_log_sync 表（跟踪已处理文件）
```

### 支持的会话来源

| 来源 | 目录 | 文件格式 |
|------|------|----------|
| Claude Code | `~/.claude/projects/` | JSONL |
| Codex | `~/.codex/sessions/` | JSONL |
| Gemini CLI | `~/.gemini/tmp/` | JSON |

---

## 核心架构

### 文件结构

```
src-tauri/src/
├── services/
│   ├── session_usage.rs          # Claude 会话同步
│   ├── session_usage_codex.rs    # Codex 会话同步
│   ├── session_usage_gemini.rs   # Gemini 会话同步
│   └── usage_stats.rs            # 通用统计函数
├── proxy/
│   └── usage/
│       ├── calculator.rs          # 费用计算
│       └── parser.rs              # Token 解析
└── commands/
    └── usage.rs                  # Tauri 命令入口
```

### 核心数据结构

```rust
// 同步结果
pub struct SessionSyncResult {
    pub imported: u32,        // 新导入记录数
    pub skipped: u32,         // 跳过（已存在）记录数
    pub files_scanned: u32,   // 扫描文件数
    pub errors: Vec<String>,  // 错误列表
}

// Token 使用量
pub struct TokenUsage {
    pub input_tokens: u32,
    pub output_tokens: u32,
    pub cache_read_tokens: u32,
    pub cache_creation_tokens: u32,
    pub model: Option<String>,
    pub message_id: Option<String>,
}

// 去重键
pub struct DedupKey<'a> {
    pub app_type: &'a str,
    pub model: &'a str,
    pub input_tokens: u32,
    pub output_tokens: u32,
    pub cache_read_tokens: u32,
    pub cache_creation_tokens: u32,
    pub created_at: i64,
}
```

---

## Claude 会话日志解析

### 文件位置

```
~/.claude/projects/
├── <project-dir>/
│   ├── *.jsonl                    # 主会话文件
│   └── <session-id>/
│       └── subagents/
│           └── *.jsonl            # 子 Agent 会话
```

### JSONL 消息格式

```json
{
  "type": "assistant",
  "message": {
    "id": "msg_0123456789abcdef",
    "model": "claude-opus-4-6",
    "usage": {
      "input_tokens": 1500,
      "output_tokens": 300,
      "cache_read_input_tokens": 5000,
      "cache_creation_input_tokens": 0
    },
    "stop_reason": "end_turn"
  },
  "timestamp": "2026-04-05T12:00:00Z",
  "sessionId": "session-abcdef"
}
```

### 解析逻辑

#### 1. 文件收集

```rust
fn collect_jsonl_files(projects_dir: &Path) -> Vec<PathBuf> {
    // 扫描三层：
    // 1. projects/<project>/*.jsonl
    // 2. projects/<project>/<session>/subagents/*.jsonl
}
```

#### 2. 增量解析

```rust
// 获取文件元数据
let file_modified = metadata_modified_nanos(&metadata);

// 检查同步状态
let (last_modified, last_offset) = get_sync_state(db, &file_path_str)?;

// 文件未变化则跳过
if file_modified <= last_modified {
    return Ok((0, 0));
}

// 从上次偏移位置继续解析
for (line_offset, line) in reader.lines().enumerate() {
    if line_offset <= last_offset as usize {
        continue;
    }
    // 解析...
}
```

#### 3. 消息过滤与去重

**只保留满足以下条件的消息**：

- `type == "assistant"`
- `message.stop_reason` 存在（是完整的 API 调用，非中间状态）
- `output_tokens > 0`

**同一 message.id 的多条消息处理**：

```rust
// 优先保留有 stop_reason 的条目
if parsed.stop_reason.is_some() && existing.stop_reason.is_none() {
    replace = true;
}
// 都有/都没有时，保留 output_tokens 更大的
else if parsed.stop_reason.is_some() == existing.stop_reason.is_some() {
    replace = parsed.output_tokens > existing.output_tokens;
}
```

#### 4. Request ID 生成

```rust
let request_id = format!("session:{}", msg.message_id);
```

---

## Codex 会话日志解析

### 文件位置

```
~/.codex/
├── sessions/
│   └── YYYY/
│       └── MM/
│           └── DD/
│               └── *.jsonl
└── archived_sessions/
    └── *.jsonl
```

### JSONL 事件格式

#### session_meta（会话元数据）

```json
{
  "type": "session_meta",
  "payload": {
    "session_id": "session-0123456789"
  },
  "timestamp": "2026-04-05T12:00:00Z"
}
```

#### turn_context（模型切换）

```json
{
  "type": "turn_context",
  "payload": {
    "model": "openai/gpt-5.4-2026-03-05",
    "info": {
      "model": "openai/gpt-5.4-2026-03-05"
    }
  }
}
```

#### event_msg（Token 统计）

```json
{
  "type": "event_msg",
  "payload": {
    "type": "token_count",
    "info": {
      "total_token_usage": {
        "input_tokens": 17934,
        "cached_input_tokens": 9600,
        "output_tokens": 454
      },
      "model": "gpt-5.4"
    }
  },
  "timestamp": "2026-04-05T12:00:00Z"
}
```

### 解析逻辑

#### 1. 状态机

```rust
struct FileParseState {
    session_id: Option<String>,     // 从 session_meta 提取
    current_model: String,          // 从 turn_context 或 token_count 提取
    prev_total: Option<CumulativeTokens>,  // 前一次累计值
    event_index: u32,               // 事件计数器（用于生成 request_id）
}
```

#### 2. 模型名归一化

```rust
fn normalize_codex_model(raw: &str) -> String {
    // 1. 转小写
    let mut name = raw.to_lowercase();

    // 2. 剥离 provider 前缀
    if let Some(pos) = name.rfind('/') {
        name = name[pos + 1..].to_string();
    }

    // 3. 剥离 ISO 日期后缀 -YYYY-MM-DD
    if name.len() > 11 {
        let suffix = &name[name.len() - 11..];
        if suffix.starts_with('-') && is_date_suffix(suffix) {
            name.truncate(name.len() - 11);
        }
    }

    // 4. 剥离紧凑日期后缀 -YYYYMMDD
    if name.len() > 9 {
        let parts: Vec<&str> = name.rsplitn(2, '-').collect();
        if parts.len() == 2 && is_compact_date(parts[0]) {
            name = parts[1].to_string();
        }
    }

    name
}
```

**示例**：
- `OpenAI/GPT-5.4-2026-03-05` → `gpt-5.4`
- `GLM-4.6-20260305` → `glm-4.6`

#### 3. Delta 计算

Codex 日志记录的是**累计值**，需要计算与前一次的差值：

```rust
fn compute_delta(prev: &Option<CumulativeTokens>, current: &CumulativeTokens) -> DeltaTokens {
    match prev {
        None => DeltaTokens {
            input: current.input as u32,
            cached_input: current.cached_input as u32,
            output: current.output as u32,
        },
        Some(p) => DeltaTokens {
            input: current.input.saturating_sub(p.input) as u32,
            cached_input: current.cached_input.saturating_sub(p.cached_input) as u32,
            output: current.output.saturating_sub(p.output) as u32,
        }
    }
}
```

**防护措施**：
- `saturating_sub` 避免溢出
- `cached_input.min(input)` 钳制异常值
- 跳过 delta 全零的记录（task 边界）

#### 4. Request ID 生成

```rust
let request_id = format!("codex_session:{}:{}", session_id_str, state.event_index);
```

---

## Gemini 会话日志解析

### 文件位置

```
~/.gemini/tmp/
└── <project-hash>/
    └── chats/
        └── session-*.json
```

### JSON 会话格式

```json
{
  "sessionId": "session-0123456789abcdef",
  "messages": [
    {
      "type": "gemini",
      "id": "msg-0123456789",
      "model": "gemini-2.5-pro",
      "tokens": {
        "input": 8522,
        "output": 29,
        "cached": 3138,
        "thoughts": 405
      },
      "timestamp": "2026-04-05T12:00:00Z"
    }
  ]
}
```

### 解析逻辑

#### 1. 全量解析（非增量）

Gemini 文件是单个 JSON 对象，每次同步全量重读。

#### 2. Token 处理

```rust
// thoughts 合并到 output（思考 token 按输出计费）
let output_tokens = tokens.output + tokens.thoughts;
```

#### 3. UPSERT 语义

Gemini 会话可能更新已有记录，使用 `ON CONFLICT UPDATE`：

```sql
INSERT INTO proxy_request_logs (...)
VALUES (...)
ON CONFLICT(request_id) DO UPDATE SET
    model = excluded.model,
    input_tokens = excluded.input_tokens,
    output_tokens = excluded.output_tokens,
    cache_read_tokens = excluded.cache_read_tokens,
    input_cost_usd = excluded.input_cost_usd,
    output_cost_usd = excluded.output_cost_usd,
    cache_read_cost_usd = excluded.cache_read_cost_usd,
    cache_creation_cost_usd = excluded.cache_creation_cost_usd,
    total_cost_usd = excluded.total_cost_usd
WHERE input_tokens != excluded.input_tokens
   OR output_tokens != excluded.output_tokens
   OR cache_read_tokens != excluded.cache_read_tokens
   OR model != excluded.model
```

#### 4. Request ID 生成

```rust
let request_id = format!("gemini_session:{}:{}", session_id_str, message_id);
```

---

## 去重机制

### 避免代理/会话日志双重计数

**问题**：同一次调用可能同时存在于：
1. 代理拦截记录（`data_source = "proxy"`）
2. 会话日志导入记录（`data_source = "session_log"`）

**解决**：插入会话日志记录前，检查是否已存在相同 token 分布的代理记录。

### DedupKey

```rust
pub struct DedupKey<'a> {
    pub app_type: &'a str,           // claude/codex/gemini
    pub model: &'a str,
    pub input_tokens: u32,
    pub output_tokens: u32,
    pub cache_read_tokens: u32,
    pub cache_creation_tokens: u32,
    pub created_at: i64,
}
```

### 检查逻辑

```rust
fn should_skip_session_insert(
    conn: &Connection,
    request_id: &str,
    key: &DedupKey,
) -> Result<bool, AppError> {
    // 1. 检查是否已存在同 request_id 的记录
    let exists = conn.query_row(
        "SELECT 1 FROM proxy_request_logs WHERE request_id = ?1",
        params![request_id],
        |_| Ok(()),
    ).is_ok();

    if exists {
        return Ok(true);
    }

    // 2. 检查是否已存在相同 token 分布的代理记录
    let proxy_exists = conn.query_row(
        "SELECT 1 FROM proxy_request_logs
         WHERE app_type = ?1
           AND model = ?2
           AND input_tokens = ?3
           AND output_tokens = ?4
           AND cache_read_tokens = ?5
           AND cache_creation_tokens = ?6
           AND ABS(created_at - ?7) < 3600  -- 1 小时内
           AND (data_source != 'session_log' OR data_source IS NULL)",
        params![
            key.app_type,
            key.model,
            key.input_tokens,
            key.output_tokens,
            key.cache_read_tokens,
            key.cache_creation_tokens,
            key.created_at,
        ],
        |_| Ok(()),
    ).is_ok();

    Ok(proxy_exists)
}
```

---

## 同步状态跟踪

### session_log_sync 表

| 字段 | 类型 | 用途 |
|------|------|------|
| `file_path` | TEXT | 会话文件绝对路径（主键） |
| `last_modified` | INTEGER | 文件修改时间（纳秒级时间戳） |
| `last_line_offset` | INTEGER | Claude/Codex：已处理行数；Gemini：已处理消息数 |
| `last_synced_at` | INTEGER | 上次同步时间（秒级时间戳） |

### 辅助函数

```rust
// 获取文件修改时间（纳秒级）
fn metadata_modified_nanos(metadata: &fs::Metadata) -> i64 {
    metadata
        .modified()
        .ok()
        .and_then(|t| t.duration_since(SystemTime::UNIX_EPOCH).ok())
        .map(|d| d.as_nanos().min(i64::MAX as u128) as i64)
        .unwrap_or(0)
}

// 获取同步状态
fn get_sync_state(db: &Database, file_path: &str) -> Result<(i64, i64), AppError> {
    let conn = lock_conn!(db.conn);
    let result = conn.query_row(
        "SELECT last_modified, last_line_offset FROM session_log_sync WHERE file_path = ?1",
        rusqlite::params![file_path],
        |row| Ok((row.get::<_, i64>(0)?, row.get::<_, i64>(1)?)),
    );
    Ok(result.unwrap_or((0, 0)))
}

// 更新同步状态
fn update_sync_state(
    db: &Database,
    file_path: &str,
    last_modified: i64,
    last_offset: i64,
) -> Result<(), AppError> {
    let now = SystemTime::now()
        .duration_since(SystemTime::UNIX_EPOCH)
        .map(|d| d.as_secs() as i64)
        .unwrap_or(0);

    let conn = lock_conn!(db.conn);
    conn.execute(
        "INSERT OR REPLACE INTO session_log_sync
         (file_path, last_modified, last_line_offset, last_synced_at)
         VALUES (?1, ?2, ?3, ?4)",
        rusqlite::params![file_path, last_modified, last_offset, now],
    )?;
    Ok(())
}
```

---

## 数据存储

### proxy_request_logs 表关键字段

| 字段 | 会话来源的值 |
|------|-------------|
| `request_id` | `session:<msg_id>` / `codex_session:<session>:<idx>` / `gemini_session:<session>:<msg_id>` |
| `provider_id` | `_session` / `_codex_session` / `_gemini_session` |
| `app_type` | `claude` / `codex` / `gemini` |
| `data_source` | `session_log` / `codex_session` / `gemini_session` |
| `provider_type` | `"session_log"` / `"codex_session"` / `"gemini_session"` |
| `session_id` | 从会话文件提取 |
| `latency_ms` | `0`（会话日志无此数据） |
| `first_token_ms` | NULL |
| `status_code` | `200` |
| `error_message` | NULL |
| `cost_multiplier` | `"1.0"` |
| `is_streaming` | `1` |

### 插入 SQL（Claude/Codex）

```sql
INSERT OR IGNORE INTO proxy_request_logs (
    request_id, provider_id, app_type, model, request_model,
    input_tokens, output_tokens, cache_read_tokens, cache_creation_tokens,
    input_cost_usd, output_cost_usd, cache_read_cost_usd, cache_creation_cost_usd, total_cost_usd,
    latency_ms, first_token_ms, status_code, error_message, session_id,
    provider_type, is_streaming, cost_multiplier, created_at, data_source
) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?14, ?15, ?16, ?17, ?18, ?19, ?20, ?21, ?22, ?23, ?24)
```

---

## 数据库 Schema 参考

### proxy_request_logs

```sql
CREATE TABLE proxy_request_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id TEXT UNIQUE NOT NULL,
    provider_id TEXT NOT NULL,
    app_type TEXT NOT NULL,
    model TEXT NOT NULL,
    request_model TEXT,
    input_tokens INTEGER DEFAULT 0,
    output_tokens INTEGER DEFAULT 0,
    cache_read_tokens INTEGER DEFAULT 0,
    cache_creation_tokens INTEGER DEFAULT 0,
    input_cost_usd TEXT DEFAULT '0',
    output_cost_usd TEXT DEFAULT '0',
    cache_read_cost_usd TEXT DEFAULT '0',
    cache_creation_cost_usd TEXT DEFAULT '0',
    total_cost_usd TEXT DEFAULT '0',
    latency_ms INTEGER DEFAULT 0,
    first_token_ms INTEGER,
    status_code INTEGER DEFAULT 200,
    error_message TEXT,
    session_id TEXT,
    provider_type TEXT,
    is_streaming INTEGER DEFAULT 0,
    cost_multiplier TEXT DEFAULT '1.0',
    created_at INTEGER NOT NULL,
    data_source TEXT,

    INDEX idx_proxy_request_logs_created_at (created_at),
    INDEX idx_proxy_request_logs_app_type (app_type),
    INDEX idx_proxy_request_logs_provider_id (provider_id)
);
```

### session_log_sync

```sql
CREATE TABLE session_log_sync (
    file_path TEXT PRIMARY KEY,
    last_modified INTEGER NOT NULL,
    last_line_offset INTEGER NOT NULL,
    last_synced_at INTEGER NOT NULL
);
```

### model_pricing

```sql
CREATE TABLE model_pricing (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id TEXT UNIQUE NOT NULL,
    display_name TEXT,
    input_cost_per_million TEXT NOT NULL,
    output_cost_per_million TEXT NOT NULL,
    cache_read_cost_per_million TEXT DEFAULT '0',
    cache_creation_cost_per_million TEXT DEFAULT '0',
    is_placeholder INTEGER DEFAULT 0,
    created_at INTEGER,
    updated_at INTEGER
);
```

---

## 费用计算

### 核心公式

```
input_cost = (input_tokens / 1,000,000) * input_price
output_cost = (output_tokens / 1,000,000) * output_price
cache_read_cost = (cache_read_tokens / 1,000,000) * cache_read_price
cache_creation_cost = (cache_creation_tokens / 1,000,000) * cache_creation_price

total_cost = (input_cost + output_cost + cache_read_cost + cache_creation_cost) * cost_multiplier
```

### ModelPricing 结构

```rust
pub struct ModelPricing {
    pub input_cost_per_million: Decimal,
    pub output_cost_per_million: Decimal,
    pub cache_read_cost_per_million: Decimal,
    pub cache_creation_cost_per_million: Decimal,
}
```

### 模糊匹配

找不到精确模型名时，尝试模糊匹配：
- 前缀匹配：`claude-opus` 匹配 `claude-opus-4-6`
- 包含匹配：`gpt-5.4` 匹配 `gpt-5.4-2026-03-05`

---

## 使用示例

### Tauri 命令调用

```rust
#[tauri::command]
pub fn sync_session_usage(state: AppState) -> Result<SessionSyncResult, AppError> {
    let mut result = session_usage::sync_claude_session_logs(&state.db)?;
    match session_usage_codex::sync_codex_usage(&state.db) {
        Ok(codex_result) => {
            result.imported += codex_result.imported;
            result.skipped += codex_result.skipped;
            result.files_scanned += codex_result.files_scanned;
            result.errors.extend(codex_result.errors);
        }
        Err(e) => result.errors.push(format!("Codex sync failed: {}", e)),
    }
    match session_usage_gemini::sync_gemini_usage(&state.db) {
        Ok(gemini_result) => {
            result.imported += gemini_result.imported;
            result.skipped += gemini_result.skipped;
            result.files_scanned += gemini_result.files_scanned;
            result.errors.extend(gemini_result.errors);
        }
        Err(e) => result.errors.push(format!("Gemini sync failed: {}", e)),
    }
    Ok(result)
}
```

---

## 关键设计决策

| 决策 | 原因 |
|------|------|
| 纳秒级 `last_modified` | 兼容旧的秒级数据，同时提供更高精度 |
| Claude 用 `stop_reason` 过滤 | JSONL 包含流式增量更新，只保留最终状态 |
| Codex 用 delta 计算 | Codex 记录累计值而非单次值 |
| Gemini 用 UPSERT | Gemini 文件是全量更新，可能需要更新已导入记录 |
| 代理记录优先 | 代理数据更准确（有 latency、真实 status_code 等） |

---

## 测试要点

1. **增量同步**：修改文件内容，只解析新增行
2. **去重**：同时存在代理记录和会话日志时不重复计数
3. **模型名归一化**：各种日期格式、大小写、前缀都能正确处理
4. **边界防护**：累计值减少时 delta 钳制为 0
5. **纯缓存命中**：只有 cache_read_tokens 的记录也应被导入

---

## 许可证

本文档基于 CC Switch 项目代码整理，遵循其开源许可证。
