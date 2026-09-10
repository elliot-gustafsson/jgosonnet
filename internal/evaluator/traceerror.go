package evaluator

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet/ast"
)

type TraceSignal struct {
	Str  string
	Rest Value
}

func (t *TraceSignal) Error() string {
	return "TRACE: " + t.Str
}

func (t *TraceSignal) Print(w io.Writer, n ast.Node) error {
	loc := n.Loc()
	if loc == nil {
		// Set filename to <unknown> here?
		_, err := fmt.Fprintf(w, "TRACE: %s\n", t.Str)
		return err
	}
	filename := filepath.Base(loc.FileName)
	line := loc.Begin.Line
	_, err := fmt.Fprintf(w, "TRACE: %s:%d %s\n", filename, line, t.Str)
	return err
}

type TraceError struct {
	Err    error
	Frames []Frame
}

func (t *TraceError) Error() string {
	var b strings.Builder

	b.WriteString(t.Err.Error())
	b.WriteByte('\n')

	for _, frame := range t.Frames {
		if frame.Pos == "" && frame.Name == "" {
			continue
		}
		if frame.Name != "" {
			fmt.Fprintf(&b, "\t%s\t%s\n", frame.Pos, frame.Name)
		} else {
			fmt.Fprintf(&b, "\t%s\n", frame.Pos)
		}
	}
	return b.String()
}

type Frame struct {
	Name string
	Pos  string
}

func WrapError(err error, node ast.Node) error {
	if err == nil {
		return nil
	}

	traceErr, ok := err.(*TraceError)
	if !ok {
		traceErr = &TraceError{
			Err:    err,
			Frames: make([]Frame, 0, 8),
		}
	}

	if node == nil || node.Loc() == nil {
		return traceErr
	}

	pos := node.Loc().String()
	if pos == "" {
		return traceErr
	}

	frame := Frame{Pos: pos}
	if ctx := node.Context(); ctx != nil {
		frame.Name = *ctx
	}

	// always append first frame or when new context boundary.
	numFrames := len(traceErr.Frames)
	if numFrames == 0 || traceErr.Frames[numFrames-1].Name != frame.Name {
		traceErr.Frames = append(traceErr.Frames, frame)
	}

	return traceErr
}
