package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

type customCore struct {
	zapcore.Core
}

func (c *customCore) With(fields []zapcore.Field) zapcore.Core {
	return &customCore{Core: c.Core.With(fields)}
}

func (c *customCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *customCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	// Filter the stack trace directly in Entry
	if ent.Stack != "" {
		filtered := filterStacktrace(ent.Stack)
		ent.Stack = beautifyStacktrace(filtered)
	}
	return c.Core.Write(ent, fields)
}

func tehranTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	loc, _ := time.LoadLocation("Asia/Tehran")
	enc.AppendString(t.In(loc).Format("2006-01-02 15:04:05"))
}

func filterStacktrace(stack string) string {
	lines := strings.Split(stack, "\n")
	var result []string
	for i := 0; i < len(lines)-1; i += 2 {
		funcLine := lines[i]
		fileLine := lines[i+1]

		result = append(result, funcLine+"\n"+fileLine)

		// stop once we hit a controller
		if strings.Contains(strings.ToLower(fileLine), "controller") {
			break
		}
	}

	if len(result) == 0 && len(lines) > 1 {
		// fallback: at least return the first frame
		result = append(result, lines[0]+"\n"+lines[1])
	}
	return strings.Join(result, "\n")
}

func beautifyStacktrace(stack string) string {
	lines := strings.Split(stack, "\n")
	var frames []string

	for i := 0; i < len(lines)-1; i += 2 {
		funcLine := strings.TrimSpace(lines[i])
		fileLine := strings.TrimSpace(lines[i+1])

		// Remove everything before "src/"
		if idx := strings.Index(fileLine, "src/"); idx != -1 {
			fileLine = fileLine[idx+4:] // +4 to skip "src/"
		}

		// Keep only function name (last part after dot)
		if idx := strings.LastIndex(funcLine, "."); idx != -1 {
			funcLine = funcLine[idx+1:]
		}

		frames = append(frames, fmt.Sprintf("%s@%s", funcLine, fileLine))
	}

	return strings.Join(frames, " | ")
}
