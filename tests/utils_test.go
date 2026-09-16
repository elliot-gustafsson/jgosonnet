package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/elliot-gustafsson/jgosonnet"
	"github.com/stretchr/testify/assert"
)

var standardExtVars = map[string]string{
	"stringVar": "2 + 2",
}

var standardExtCode = map[string]string{
	"codeVar":               "3 + 3",
	"errorVar":              "error 'xxx'",
	"staticErrorVar":        ")",
	"UndeclaredX":           "x",
	"selfRecursiveVar":      `[42, std.extVar("selfRecursiveVar")[0] + 1]`,
	"mutuallyRecursiveVar1": `[42, std.extVar("mutuallyRecursiveVar2")[0] + 1]`,
	"mutuallyRecursiveVar2": `[42, std.extVar("mutuallyRecursiveVar1")[0] + 1]`,
}

func mustReadFile(t *testing.T, relPath string) string {
	t.Helper()
	b, err := os.ReadFile(relPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", relPath, err)
	}
	return string(b)
}

func readGoldenDir(t *testing.T, relDir string) string {
	t.Helper()
	entries, err := os.ReadDir(relDir)
	if err != nil {
		t.Fatalf("failed to read golden dir %s: %v", relDir, err)
	}
	m := make(map[string]string)
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join(relDir, e.Name()))
		if err != nil {
			t.Fatalf("failed to read file %s: %v", e.Name(), err)
		}
		m[e.Name()] = string(b)
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	return string(data)
}

type TestConfig struct {
	WorkDir      string
	File         string
	Snippet      string
	Expected     string
	ExpectedErr  string
	ExtVars      map[string]string
	ExtCodes     map[string]string
	TLAVars      map[string]string
	TLACodes     map[string]string
	IsMulti      bool
	NoNewline    bool
	StringOutput bool
	MaxStack     uint32
}

func runTest(t *testing.T, cfg TestConfig) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("evaluation panicked: %v", r)
		}
	}()

	if cfg.Expected != "" && cfg.ExpectedErr != "" {
		t.Fatalf("both Expected and ExpectedErr are set in TestConfig for %s", cfg.File)
	}

	origWd, _ := os.Getwd()
	absWorkdir, err := filepath.Abs(cfg.WorkDir)
	if err != nil {
		t.Fatalf("abs workdir error: %v", err)
	}
	if err := os.Chdir(absWorkdir); err != nil {
		t.Fatalf("chdir error: %v", err)
	}
	defer os.Chdir(origWd)

	evalPath := cfg.File

	type result struct {
		stdout string
		multi  map[string]string
		err    error
	}
	ch := make(chan result, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- result{err: fmt.Errorf("panic: %v", r)}
			}
		}()

		ev := jgosonnet.NewEvaluator()

		if cfg.MaxStack > 0 {
			ev.MaxStack(cfg.MaxStack)
		}
		if cfg.NoNewline {
			ev.NoNewline(true)
		}
		for k, v := range cfg.ExtVars {
			ev.ExtVar(k, v)
		}
		for k, v := range cfg.ExtCodes {
			ev.ExtCode(k, v)
		}
		for k, v := range cfg.TLAVars {
			ev.TLAVar(k, v)
		}
		for k, v := range cfg.TLACodes {
			ev.TLACode(k, v)
		}

		if cfg.IsMulti {
			var res map[string]string
			var err error
			if cfg.StringOutput {
				res, err = ev.EvaluateStringMulti(evalPath)
			} else {
				res, err = ev.EvaluateJsonMulti(evalPath)
			}
			ch <- result{multi: res, err: err}
		} else {
			var traceBuf bytes.Buffer
			ev.TraceOut(&traceBuf)
			var res string
			var err error
			if cfg.StringOutput {
				res, err = ev.EvaluateString(evalPath)
			} else {
				res, err = ev.EvaluateJson(evalPath)
			}
			ch <- result{stdout: traceBuf.String() + res, err: err}
		}
	}()

	var res result
	select {
	case res = <-ch:
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("evaluation timed out (infinite recursion/loop)")
	}

	if cfg.ExpectedErr != "" {
		if res.err == nil {
			t.Fatalf("expected error but evaluation succeeded with output:\n%s", res.stdout)
		}
		assert.Equal(t, cfg.ExpectedErr, res.err.Error())
		return
	}

	if res.err != nil {
		t.Fatalf("evaluation failed:\n%v", res.err)
	}

	if cfg.IsMulti {
		var expectedFiles map[string]string
		if err := json.Unmarshal([]byte(cfg.Expected), &expectedFiles); err != nil {
			t.Fatalf("unmarshal expected multi output error: %v", err)
		}
		for k, expContent := range expectedFiles {
			actContent, ok := res.multi[k]
			if !ok {
				t.Errorf("missing output file: %s", k)
				continue
			}
			assert.Equal(t, expContent, actContent)
		}
		return
	}

	assert.Equal(t, cfg.Expected, res.stdout)
}
