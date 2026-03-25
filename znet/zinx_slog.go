// Package znet provides the ZinxHandler, a custom slog.Handler for the zinx framework.
// It formats log output in zinx's traditional style and supports file rotation via zutils.Writer.
package znet

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aceld/zinx/v3/zutils"
)

// ZinxHandler is a custom slog.Handler that formats logs in zinx's traditional style.
// Output format: 2006/01/02 15:04:05 [LEVEL] file.go:123: message key=value
type ZinxHandler struct {
	mu        sync.Mutex
	writer    io.Writer
	opts      *slog.HandlerOptions
	prefix    string
	attrs     []slog.Attr
	groups    []string
	showDate  bool
	showTime  bool
	showFile  bool
	showLevel bool
	fw        *zutils.Writer
}

// ZinxHandlerOptions configures a ZinxHandler.
type ZinxHandlerOptions struct {
	// Level sets the minimum log level. Default is slog.LevelDebug.
	Level slog.Leveler
	// Prefix is prepended to each log line wrapped in angle brackets.
	Prefix string
	// ShowDate enables date output (2006/01/02). Default true.
	ShowDate bool
	// ShowTime enables time output (15:04:05). Default true.
	ShowTime bool
	// ShowFile enables file:line output. Default true.
	ShowFile bool
	// ShowLevel enables log level output ([INFO] etc). Default true.
	ShowLevel bool
}

// NewZinxHandler creates a new ZinxHandler writing to w.
// Pass nil opts to use defaults (all fields shown, LevelDebug).
func NewZinxHandler(w io.Writer, opts *ZinxHandlerOptions) *ZinxHandler {
	if opts == nil {
		opts = &ZinxHandlerOptions{
			ShowDate:  true,
			ShowTime:  true,
			ShowFile:  true,
			ShowLevel: true,
		}
	}
	level := slog.LevelDebug
	if opts.Level != nil {
		level = opts.Level.Level()
	}
	if w == nil {
		w = os.Stderr
	}
	return &ZinxHandler{
		writer: w,
		opts: &slog.HandlerOptions{
			Level: level,
		},
		prefix:    opts.Prefix,
		showDate:  opts.ShowDate,
		showTime:  opts.ShowTime,
		showFile:  opts.ShowFile,
		showLevel: opts.ShowLevel,
	}
}

// NewZinxFileHandler creates a ZinxHandler that writes to a rotating file.
func NewZinxFileHandler(fileDir, fileName string, opts *ZinxHandlerOptions) *ZinxHandler {
	h := NewZinxHandler(os.Stderr, opts)
	h.SetLogFile(fileDir, fileName)
	return h
}

// SetLogFile redirects output to a rotating log file.
func (h *ZinxHandler) SetLogFile(fileDir, fileName string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.fw != nil {
		h.fw.Close()
	}
	h.fw = zutils.New(filepath.Join(fileDir, fileName))
	h.writer = h.fw
}

// SetMaxAge sets maximum retention days for rotated log files.
func (h *ZinxHandler) SetMaxAge(ma int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.fw != nil {
		h.fw.SetMaxAge(ma)
	}
}

// SetMaxSize sets maximum size in bytes per log file before rotation.
func (h *ZinxHandler) SetMaxSize(ms int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.fw != nil {
		h.fw.SetMaxSize(ms)
	}
}

// SetCons enables simultaneous output to stderr when writing to a file.
func (h *ZinxHandler) SetCons(b bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.fw != nil {
		h.fw.SetCons(b)
	}
}

// SetLevel sets the minimum log level dynamically.
func (h *ZinxHandler) SetLevel(level slog.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.opts == nil {
		h.opts = &slog.HandlerOptions{}
	}
	h.opts.Level = level
}

// SetPrefix sets the log prefix shown as <prefix> at the start of each line.
func (h *ZinxHandler) SetPrefix(prefix string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.prefix = prefix
}

// Enabled reports whether the handler handles records at the given level.
func (h *ZinxHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelDebug
	if h.opts != nil && h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats and writes the log record.
func (h *ZinxHandler) Handle(ctx context.Context, r slog.Record) error {
	var file string
	var line int
	if h.showFile && r.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := frames.Next()
		file = f.File
		line = f.Line
	}

	buf := make([]byte, 0, 256)

	if h.prefix != "" {
		buf = append(buf, '<')
		buf = append(buf, h.prefix...)
		buf = append(buf, '>')
		buf = append(buf, ' ')
	}

	if h.showDate {
		year, month, day := r.Time.Date()
		buf = zinxAppendInt(buf, year, 4)
		buf = append(buf, '/')
		buf = zinxAppendInt(buf, int(month), 2)
		buf = append(buf, '/')
		buf = zinxAppendInt(buf, day, 2)
		buf = append(buf, ' ')
	}

	if h.showTime {
		hour, min, sec := r.Time.Clock()
		buf = zinxAppendInt(buf, hour, 2)
		buf = append(buf, ':')
		buf = zinxAppendInt(buf, min, 2)
		buf = append(buf, ':')
		buf = zinxAppendInt(buf, sec, 2)
		buf = append(buf, ' ')
	}

	if h.showLevel {
		buf = append(buf, '[')
		buf = append(buf, zinxLevelString(r.Level)...)
		buf = append(buf, ']')
		buf = append(buf, ' ')
	}

	if h.showFile && file != "" {
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
		buf = append(buf, ':')
		buf = append(buf, ' ')
	}

	buf = append(buf, r.Message...)

	for _, a := range h.attrs {
		buf = zinxAppendAttr(buf, a, h.groups)
	}
	r.Attrs(func(a slog.Attr) bool {
		buf = zinxAppendAttr(buf, a, h.groups)
		return true
	})

	buf = zinxAppendContextFields(ctx, buf)
	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	w := h.writer
	if w == nil {
		w = os.Stderr
	}
	_, err := w.Write(buf)
	return err
}

// WithAttrs returns a new handler with the given attributes pre-set.
func (h *ZinxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := h.clone()
	h2.attrs = make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(h2.attrs, h.attrs)
	copy(h2.attrs[len(h.attrs):], attrs)
	return h2
}

// WithGroup returns a new handler scoped to the given group name.
func (h *ZinxHandler) WithGroup(name string) slog.Handler {
	h2 := h.clone()
	h2.groups = make([]string, len(h.groups)+1)
	copy(h2.groups, h.groups)
	h2.groups[len(h.groups)] = name
	return h2
}

func (h *ZinxHandler) clone() *ZinxHandler {
	return &ZinxHandler{
		writer:    h.writer,
		opts:      h.opts,
		prefix:    h.prefix,
		attrs:     h.attrs,
		groups:    h.groups,
		showDate:  h.showDate,
		showTime:  h.showTime,
		showFile:  h.showFile,
		showLevel: h.showLevel,
		fw:        h.fw,
	}
}

// defaultZinxHandler is the global zinx slog handler.
var defaultZinxHandler = NewZinxHandler(os.Stderr, &ZinxHandlerOptions{
	ShowDate:  true,
	ShowTime:  true,
	ShowFile:  true,
	ShowLevel: true,
})

func init() {
	slog.SetDefault(slog.New(defaultZinxHandler))
}

// SetupSlog replaces the global slog default with a custom handler.
func SetupSlog(handler slog.Handler) {
	slog.SetDefault(slog.New(handler))
}

// SetLogFile configures file output on the default handler.
func SetLogFile(fileDir, fileName string) {
	defaultZinxHandler.SetLogFile(fileDir, fileName)
}

// SetMaxAge sets maximum retention days on the default handler.
func SetMaxAge(ma int) {
	defaultZinxHandler.SetMaxAge(ma)
}

// SetMaxSize sets maximum file size (bytes) on the default handler.
func SetMaxSize(ms int64) {
	defaultZinxHandler.SetMaxSize(ms)
}

// SetCons enables console output on the default handler.
func SetCons(b bool) {
	defaultZinxHandler.SetCons(b)
}

// SetLogLevel sets the minimum log level on the default handler.
// Accepts slog.Level values directly.
func SetLogLevel(level slog.Level) {
	defaultZinxHandler.SetLevel(level)
}

// SetLogPrefix sets the log prefix on the default handler.
func SetLogPrefix(prefix string) {
	defaultZinxHandler.SetPrefix(prefix)
}

// zinxLog writes a log record with the caller at the given skip depth above zinxLog itself.
// skip=1 means the direct caller of zinxLog appears in the log.
func zinxLog(skip int, level slog.Level, msg string) {
	logger := slog.Default()
	if !logger.Enabled(nil, level) {
		return
	}
	var pcs [1]uintptr
	// skip: runtime.Callers(1) + zinxLog(1) + skip
	runtime.Callers(skip+2, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	_ = logger.Handler().Handle(nil, r)
}

// ContextKey is the type for context keys used in zinx logging.
type ContextKey string

const (
	ContextKeyTraceID   ContextKey = "trace_id"
	ContextKeySpanID    ContextKey = "span_id"
	ContextKeyRequestID ContextKey = "request_id"
	ContextKeyUserID    ContextKey = "user_id"
)

func zinxLevelString(l slog.Level) string {
	switch {
	case l < slog.LevelInfo:
		return "DEBUG"
	case l < slog.LevelWarn:
		return "INFO"
	case l < slog.LevelError:
		return "WARN"
	default:
		return "ERROR"
	}
}

func zinxAppendInt(buf []byte, i, width int) []byte {
	s := strconv.Itoa(i)
	for len(s) < width {
		s = "0" + s
	}
	return append(buf, s...)
}

func zinxAppendAttr(buf []byte, a slog.Attr, groups []string) []byte {
	if a.Equal(slog.Attr{}) {
		return buf
	}
	buf = append(buf, ' ')
	if len(groups) > 0 {
		buf = append(buf, strings.Join(groups, ".")...)
		buf = append(buf, '.')
	}
	buf = append(buf, a.Key...)
	buf = append(buf, '=')
	buf = append(buf, fmt.Sprintf("%v", a.Value.Any())...)
	return buf
}

func zinxAppendContextFields(ctx context.Context, buf []byte) []byte {
	if ctx == nil {
		return buf
	}
	if v, ok := ctx.Value(ContextKeyTraceID).(string); ok && v != "" {
		buf = append(buf, " trace_id="...)
		buf = append(buf, v...)
	}
	if v, ok := ctx.Value(ContextKeySpanID).(string); ok && v != "" {
		buf = append(buf, " span_id="...)
		buf = append(buf, v...)
	}
	if v, ok := ctx.Value(ContextKeyRequestID).(string); ok && v != "" {
		buf = append(buf, " request_id="...)
		buf = append(buf, v...)
	}
	if v, ok := ctx.Value(ContextKeyUserID).(string); ok && v != "" {
		buf = append(buf, " user_id="...)
		buf = append(buf, v...)
	}
	return buf
}
