# Zlog 迁移至 Slog 升级方案

## 1. 现状分析

### 1.1 项目结构概述
- 项目：Zinx (github.com/aceld/zinx/v3) - 基于 Golang 的轻量级并发服务器框架
- Go 版本：1.26.1
- 当前日志库：项目内部自定义的 `zlog` 包（位于 `zlog/` 目录）

### 1.2 Zlog 依赖引用
#### 导入 `github.com/aceld/zinx/v3/zlog` 的文件（共 54 个）：
- 核心网络组件：`znet/client.go`, `znet/connection.go`, `znet/kcp_connection.go`, `znet/ws_connection.go`, `znet/server.go`, `znet/connmanager.go`, `znet/msghandler.go`, `znet/heartbeat.go`, `znet/middleware.go`, `znet/defaultrouterfunc.go`
- 其他组件：`znotify/notify.go`, `ztimer/delayfunc.go`, `ztimer/timerscheduler.go`, `ztimer/timewheel.go`, `zdecoder/htlvcrcdecoder.go`, `zasync_op/async_worker.go`, `zconf/userconf.go`, `zconf/zconf.go`
- 示例代码：`examples/` 目录下的多个示例程序（约 20 个文件）
- 测试文件：`zlog/zlog_test.go`, `ztimer/timerscheduler_test.go`

### 1.3 Zlog 包核心组件
#### 1.3.1 `zlog/logger_core.go`
- 定义了 `ZinxLoggerCore` 结构体，支持多协程安全的日志记录
- 提供日志级别：Debug, Info, Warn, Error, Panic, Fatal
- 支持日志头部格式自定义（日期、时间、文件名、行号等）
- 支持日志文件输出、轮转（通过 `zutils.Writer`）
- 支持日志钩子（hook）功能

#### 1.3.2 `zlog/stdzlog.go`
- 提供全局日志实例 `StdZinxLog`
- 提供全局日志方法（Debug, Info, Warn, Error, Fatal, Panic, Stack）
- 支持日志级别、前缀、文件输出等配置

#### 1.3.3 `zlog/default.go`
- 实现了 `ziface.ILogger` 接口的默认日志器
- 提供 `InfoF`, `ErrorF`, `DebugF` 方法（无 context）
- 提供 `InfoFX`, `ErrorFX`, `DebugFX` 方法（带 context，但仅打印 context 未实际集成）
- 提供 `SetLogger` 和 `Ins()` 全局访问方法

### 1.4 ILogger 接口 (`ziface/ilogger.go`)
```go
type ILogger interface {
    //without context
    InfoF(format string, v ...interface{})
    ErrorF(format string, v ...interface{})
    DebugF(format string, v ...interface{})
    
    //with context
    InfoFX(ctx context.Context, format string, v ...interface{})
    ErrorFX(ctx context.Context, format string, v ...interface{})
    DebugFX(ctx context.Context, format string, v ...interface{})
}
```

### 1.5 现有日志调用模式
1. **全局函数调用**（最常见）：
   ```go
   zlog.Debug("message")
   zlog.Infof("format %s", arg)
   zlog.Error("error occurred")
   ```

2. **实例方法调用**：
   ```go
   zlog.Ins().InfoF("message")
   zlog.Ins().ErrorF("error: %v", err)
   ```

3. **带 Context 的调用**（较少使用）：
   ```go
   zlog.Ins().InfoFX(ctx, "message")
   ```

### 1.6 现有 Slog 支持
- 已在 `znet/middleware.go` 中实现 `SlogLoggerMiddleware` 和 `SlogLoggerMiddlewareWithLevel`
- 支持 OpenTelemetry 集成（`OTelTraceMiddleware`）
- 支持 TraceID 追踪（通过 context 或 c.Set("traceID", ...)）

### 1.7 日志文件轮转
- 使用 `zutils.Writer` 实现日志文件轮转（按日期和大小）
- 支持最大保留天数、单个日志最大容量、同时输出控制台等功能

## 2. 迁移目标

### 2.1 核心目标
1. **彻底移除 `zlog` 包**，替换为 Go 标准库 `log/slog`
2. **实现与 `context.Context` 的深度集成**，支持日志记录中的 TraceID 追踪及其他上下文关联字段
3. **保持现有日志功能**：日志级别、文件输出、轮转、格式化等
4. **最小化代码变更**，尽量通过适配器模式减少对业务代码的影响

### 2.2 具体目标
1. **依赖清理**：
   - 移除 `github.com/aceld/zinx/v3/zlog` 依赖
   - 保留 `zutils.Writer` 用于日志文件轮转
   - 可能需要调整 `ziface.ILogger` 接口以适应 slog

2. **Slog 初始化配置**：
   - 设计自定义 slog.Handler，支持现有日志格式
   - 支持文件输出和轮转（通过 `zutils.Writer`）
   - 支持日志级别动态调整
   - 支持 TraceID 等上下文字段自动提取

3. **替换策略**：
   - 全局日志实例替换为 slog 默认 logger
   - 提供适配器函数，将现有 `zlog.Debug()` 等调用转换为 slog 调用
   - 逐步迁移，确保兼容性

4. **Context 集成**：
   - 在中间件中自动提取 TraceID 并添加到 slog logger
   - 支持通过 context 传递自定义字段
   - 确保日志方法能正确接收 context 参数

## 3. 详细规划

### 3.1 阶段一：准备工作（1-2天）
#### 3.1.1 创建 Slog 自定义 Handler
- 设计 `ZinxSlogHandler` 结构体，嵌入 `slog.TextHandler` 或自定义实现
- 支持以下特性：
  - 自定义时间格式（与现有格式兼容）
  - 自定义日志级别显示（如 `[DEBUG]`, `[INFO]` 等）
  - 支持前缀字段
  - 支持文件名和行号（通过 `runtime.Caller`）
  - 支持通过 context 提取 TraceID 等字段
  - 支持输出到 `zutils.Writer` 实现文件轮转

#### 3.1.2 实现 Slog Handler 核心方法
```go
type ZinxSlogHandler struct {
    opts   *slog.HandlerOptions
    prefix string
    writer *zutils.Writer
    // 其他配置字段
}

func (h *ZinxSlogHandler) Enabled(ctx context.Context, level slog.Level) bool
func (h *ZinxSlogHandler) Handle(ctx context.Context, r slog.Record) error
func (h *ZinxSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler
func (h *ZinxSlogHandler) WithGroup(name string) slog.Handler
```

#### 3.1.3 设计 Context 字段提取机制
- 定义 context key 类型（避免冲突）
- 实现从 context 中提取 TraceID、SpanID 等字段
- 支持 OpenTelemetry trace context 提取

### 3.2 阶段二：适配器层实现（2-3天）
#### 3.1.1 创建 Zlog 适配器
- 在 `zlog` 包中创建适配器，将现有 API 转换为 slog 调用
- 保持 `zlog.Debug()`, `zlog.Infof()` 等函数签名不变
- 内部实现转换为 slog 调用

#### 3.2.2 修改 `zlog/default.go`
- 更新 `zinxDefaultLog` 实现，内部使用 slog
- 保持 `ILogger` 接口不变，但实现改为 slog
- 确保 `InfoFX` 等方法能正确传递 context

#### 3.2.3 更新 `zlog/stdzlog.go`
- 将全局 `StdZinxLog` 实例改为基于 slog 的实现
- 保持全局函数接口不变

### 3.3 阶段三：逐步迁移（3-5天）
#### 3.3.1 迁移核心组件（优先级高）
1. `znet/server.go` - 服务器启动日志
2. `znet/connection.go` - 连接管理日志
3. `znet/client.go` - 客户端日志
4. `znet/middleware.go` - 中间件日志（已有 slog 支持，需调整）

#### 3.3.2 迁移其他组件
1. `znet/connmanager.go`
2. `znet/msghandler.go`
3. `znet/heartbeat.go`
4. `znotify/notify.go`
5. `ztimer/` 包相关文件
6. `zdecoder/` 包相关文件
7. `zasync_op/` 包相关文件
8. `zconf/` 包相关文件

#### 3.3.3 迁移示例代码
- 更新 `examples/` 目录下的所有示例程序
- 确保示例代码仍能正常工作

### 3.4 阶段四：清理与优化（1-2天）
#### 3.4.1 移除 Zlog 包
- 删除 `zlog/` 目录（或保留为空壳适配器）
- 更新 `go.mod` 移除相关依赖
- 更新文档和注释

#### 3.4.2 性能优化
- 基准测试，确保 slog 性能满足要求
- 优化日志格式化性能
- 考虑异步日志记录（如需要）

#### 3.4.3 功能增强
- 支持结构化日志（JSON 格式）
- 支持动态日志级别调整
- 支持日志采样（如需要）

### 3.5 具体实现细节

#### 3.5.1 Slog Handler 设计要点
```go
// 示例实现框架
type ZinxHandler struct {
    mu        sync.Mutex
    writer    io.Writer
    opts      slog.HandlerOptions
    prefix    string
    callDepth int
}

func NewZinxHandler(w io.Writer, opts *slog.HandlerOptions) *ZinxHandler {
    // 初始化
}

func (h *ZinxHandler) Handle(ctx context.Context, r slog.Record) error {
    // 1. 获取 caller 信息
    // 2. 格式化时间
    // 3. 提取 context 字段（TraceID 等）
    // 4. 格式化日志行（参考现有格式）
    // 5. 写入 writer
}
```

#### 3.5.2 Context 集成方案
```go
// 定义 context key
type contextKey string

const (
    TraceIDKey contextKey = "trace_id"
    SpanIDKey  contextKey = "span_id"
)

// 从 context 提取字段
func extractContextFields(ctx context.Context) []slog.Attr {
    var attrs []slog.Attr
    if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
        attrs = append(attrs, slog.String("trace_id", traceID))
    }
    // 其他字段...
    return attrs
}
```

#### 3.5.3 文件轮转集成
- 创建 `ZinxWriterHandler`，包装 `zutils.Writer`
- 支持现有配置：`SetLogFile`, `SetMaxAge`, `SetMaxSize`, `SetCons`
- 保持向后兼容

### 3.6 兼容性测试方案

#### 3.6.1 单元测试
1. **Handler 测试**：
   - 测试日志级别过滤
   - 测试日志格式输出
   - 测试 context 字段提取
   - 测试文件输出

2. **适配器测试**：
   - 测试 `zlog.Debug()` -> slog 转换
   - 测试 `zlog.Ins().InfoF()` -> slog 转换
   - 测试带 context 的调用

3. **集成测试**：
   - 测试中间件日志记录
   - 测试 TraceID 传递
   - 测试日志文件轮转

#### 3.6.2 性能测试
1. **基准测试**：
   - 对比 slog 与现有 zlog 性能
   - 测试不同日志级别的性能影响
   - 测试文件输出性能

2. **压力测试**：
   - 高并发日志记录测试
   - 大日志文件轮转测试

#### 3.6.3 兼容性验证
1. **格式兼容性**：
   - 确保日志格式与现有格式一致（或明确变更）
   - 确保日志解析工具仍能工作

2. **功能兼容性**：
   - 确保所有现有日志功能正常
   - 确保配置方法向后兼容

### 3.7 回滚方案
- 保留 `zlog` 包的 git 历史
- 准备回滚脚本，必要时恢复旧代码
- 采用渐进式迁移，降低风险

## 4. 输出文档
本文档已保存至 `upgrade_log.md`，包含完整的分析结果和详细的升级步骤。

### 后续步骤：
1. 评审本方案，确认迁移范围和优先级
2. 创建分支开始实施
3. 按阶段执行迁移，每个阶段完成后进行测试
4. 完成后更新相关文档

## 附录 A：自定义 Slog Handler 设计示例

### A.1 设计目标
创建一个与现有 zlog 格式兼容的 slog.Handler，支持：
- 自定义时间格式（`2006/01/02 15:04:05`）
- 日志级别显示为 `[DEBUG]`、`[INFO]` 等
- 支持前缀字段（如 `<PREFIX>`）
- 自动提取文件名和行号（通过 `runtime.Caller`）
- 支持从 context 提取 TraceID 等字段
- 支持输出到 `zutils.Writer` 实现文件轮转

### A.2 核心结构体
```go
// zlog/zinx_handler.go
package zlog

import (
    "context"
    "io"
    "log/slog"
    "runtime"
    "strconv"
    "sync"
    "time"
)

type ZinxHandler struct {
    mu        sync.Mutex
    writer    io.Writer
    opts      *slog.HandlerOptions
    prefix    string
    callDepth int
    attrs     []slog.Attr
    groups    []string
}

func NewZinxHandler(w io.Writer, opts *slog.HandlerOptions) *ZinxHandler {
    if opts == nil {
        opts = &slog.HandlerOptions{}
    }
    return &ZinxHandler{
        writer:    w,
        opts:      opts,
        callDepth: 4, // 默认调用深度
    }
}
```

### A.3 核心方法实现
```go
func (h *ZinxHandler) Enabled(ctx context.Context, level slog.Level) bool {
    return level >= h.opts.Level.Level()
}

func (h *ZinxHandler) Handle(ctx context.Context, r slog.Record) error {
    // 1. 获取 caller 信息
    var file string
    var line int
    if r.PC != 0 {
        fs := runtime.CallersFrames([]uintptr{r.PC})
        f, _ := fs.Next()
        file = f.File
        line = f.Line
    }
    
    // 2. 格式化时间
    timeStr := r.Time.Format("2006/01/02 15:04:05")
    
    // 3. 提取 context 字段（TraceID 等）
    var contextAttrs []slog.Attr
    if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
        contextAttrs = append(contextAttrs, slog.String("trace_id", traceID))
    }
    
    // 4. 构建日志行
    buf := make([]byte, 0, 256)
    
    // 前缀
    if h.prefix != "" {
        buf = append(buf, '<')
        buf = append(buf, h.prefix...)
        buf = append(buf, '>')
    }
    
    // 时间
    buf = append(buf, timeStr...)
    buf = append(buf, ' ')
    
    // 级别
    buf = append(buf, '[')
    buf = append(buf, r.Level.String()...)
    buf = append(buf, ' ')
    
    // 文件名和行号
    if file != "" {
        short := file
        for i := len(file) - 1; i > 0; i-- {
            if file[i] == '/' {
                short = file[i+1:]
                break
            }
        }
        buf = append(buf, short...)
        buf = append(buf, ':')
        buf = append(buf, strconv.Itoa(line)...)
        buf = append(buf, ' ')
    }
    
    // 消息
    buf = append(buf, r.Message...)
    
    // 属性
    r.Attrs(func(a slog.Attr) bool {
        buf = append(buf, ' ')
        buf = append(buf, a.Key...)
        buf = append(buf, '=')
        buf = append(buf, a.Value.String()...)
        return true
    })
    
    // context 属性
    for _, a := range contextAttrs {
        buf = append(buf, ' ')
        buf = append(buf, a.Key...)
        buf = append(buf, '=')
        buf = append(buf, a.Value.String()...)
    }
    
    buf = append(buf, '\n')
    
    // 5. 写入 writer
    h.mu.Lock()
    defer h.mu.Unlock()
    _, err := h.writer.Write(buf)
    return err
}

func (h *ZinxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    h2 := *h
    h2.attrs = make([]slog.Attr, len(h.attrs)+len(attrs))
    copy(h2.attrs, h.attrs)
    copy(h2.attrs[len(h.attrs):], attrs)
    return &h2
}

func (h *ZinxHandler) WithGroup(name string) slog.Handler {
    h2 := *h
    h2.groups = make([]string, len(h.groups)+1)
    copy(h2.groups, h.groups)
    h2.groups[len(h.groups)] = name
    return &h2
}
```

### A.4 使用示例
```go
// 初始化
writer := zutils.New("./logs/zinx.log")
handler := NewZinxHandler(writer, &slog.HandlerOptions{
    Level: slog.LevelDebug,
})
logger := slog.New(handler)
slog.SetDefault(logger)

// 使用
slog.Info("Server started", "port", 8999)
slog.Debug("Debug message", "key", "value")

// 带 context
ctx := context.WithValue(context.Background(), TraceIDKey, "123456")
slog.InfoContext(ctx, "Request processed")
```

### A.5 与现有 zlog 配置的映射
| zlog 配置 | slog Handler 实现 |
|-----------|-------------------|
| `SetLogFile(fileDir, fileName)` | 创建 `zutils.Writer` 并传递给 Handler |
| `SetMaxAge(ma)` | 调用 `writer.SetMaxAge(ma)` |
| `SetMaxSize(ms)` | 调用 `writer.SetMaxSize(ms)` |
| `SetCons(b)` | 调用 `writer.SetCons(b)` |
| `SetPrefix(prefix)` | Handler 中设置 prefix 字段 |
| `ResetFlags(flag)` | Handler 中配置时间、文件等显示选项 |
| `SetLogLevel(level)` | 通过 `slog.HandlerOptions.Level` 设置 |

## 附录 B：Context 集成详细设计

### B.1 Context Key 定义
```go
// 在 ziface 包中定义
type ContextKey string

const (
    TraceIDKey ContextKey = "trace_id"
    SpanIDKey  ContextKey = "span_id"
    RequestIDKey ContextKey = "request_id"
    UserIDKey  ContextKey = "user_id"
)
```

### B.2 中间件集成
```go
// 在中间件中自动提取 TraceID
func SlogLoggerMiddleware() ziface.HandlerFunc {
    return func(c *ziface.Context) {
        // 从 OpenTelemetry span 中提取 trace ID
        if span := trace.SpanFromContext(c.Ctx); span.SpanContext().HasTraceID() {
            traceID := span.SpanContext().TraceID().String()
            // 添加到 context
            c.Ctx = context.WithValue(c.Ctx, ziface.TraceIDKey, traceID)
        }
        
        // 创建带有基础字段的 logger
        logger := slog.Default().With(
            "msgID", c.MsgID,
            "connID", c.Conn.GetConnID(),
        )
        
        // 将 logger 存储到 context 中
        c.Ctx = context.WithValue(c.Ctx, "logger", logger)
        
        // 继续处理
        c.Next()
    }
}
```

### B.3 日志记录示例
```go
// 在处理函数中使用
func HandleMessage(c *ziface.Context) {
    // 从 context 获取 logger
    if logger, ok := c.Ctx.Value("logger").(*slog.Logger); ok {
        logger.InfoContext(c.Ctx, "Processing message")
    }
    
    // 或者直接使用 slog 全局 logger
    slog.InfoContext(c.Ctx, "Processing message")
}
```

---
*文档生成时间：2026-03-23*
*分析工具：OpenCode AI Assistant*