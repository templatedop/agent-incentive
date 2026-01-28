package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// helper to extract caller from a single-line JSON log
func extractCaller(t *testing.T, s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) == 0 {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &m); err != nil { // last line in case of extra newlines
		// Fail the test to surface log format changes early
		// Use t.Fatalf for immediate stop
		// but still return something for static analysis
		// NOTE: if format changes, adjust this test.
		return ""
	}
	caller, _ := m["caller"].(string)
	return caller
}

func TestCallerSimpleAPI(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf)
	baseLogger = &Logger{logger: &zl} // override global for test

	// Determine expected file & line: next line will call Info()
	_, expectedFile, expectedLineBefore, _ := runtime.Caller(0)
	expectedLine := expectedLineBefore + 1
	Info(nil, "test simple caller") // THIS line must match expectedLine

	caller := extractCaller(t, buf.String())
	if caller == "" {
		t.Fatalf("caller field missing in log output: %s", buf.String())
	}
	// caller format is path:line
	if !strings.Contains(caller, ":") {
		t.Fatalf("unexpected caller format: %s", caller)
	}
	parts := strings.Split(caller, ":")
	lineStr := parts[len(parts)-1]
	if lineStr != intToString(expectedLine) || !strings.HasSuffix(parts[0], filepathBase(expectedFile)) {
		// Provide detailed diff
		if lineStr != intToString(expectedLine) {
			t.Errorf("line mismatch: got %s want %d", lineStr, expectedLine)
		}
		if !strings.HasSuffix(parts[0], filepathBase(expectedFile)) {
			t.Errorf("file mismatch: got %s want suffix %s", parts[0], filepathBase(expectedFile))
		}
	}
}

func TestCallerEventAPI(t *testing.T) {
	var buf bytes.Buffer
	zl := zerolog.New(&buf)
	baseLogger = &Logger{logger: &zl}

	_, expectedFile, expectedLineBefore, _ := runtime.Caller(0)
	expectedLine := expectedLineBefore + 1
	InfoEvent(nil).Msg("test event caller") // THIS line must match expectedLine

	caller := extractCaller(t, buf.String())
	if caller == "" {
		t.Fatalf("caller field missing in log output: %s", buf.String())
	}
	parts := strings.Split(caller, ":")
	lineStr := parts[len(parts)-1]
	if lineStr != intToString(expectedLine) || !strings.HasSuffix(parts[0], filepathBase(expectedFile)) {
		if lineStr != intToString(expectedLine) {
			t.Errorf("line mismatch: got %s want %d", lineStr, expectedLine)
		}
		if !strings.HasSuffix(parts[0], filepathBase(expectedFile)) {
			t.Errorf("file mismatch: got %s want suffix %s", parts[0], filepathBase(expectedFile))
		}
	}
}

// intToString avoids importing strconv solely for tests (minimal footprint)
func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}

// filepathBase avoids importing path/filepath; we can use strings since we only need suffix.
func filepathBase(p string) string {
	// Normalize separators: windows may use '\\'
	p = strings.ReplaceAll(p, "\\", "/")
	idx := strings.LastIndex(p, "/")
	if idx == -1 {
		return p
	}
	return p[idx+1:]
}
