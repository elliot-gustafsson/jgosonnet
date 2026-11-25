package evaluator

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"unsafe"

	fast "github.com/elliot-gustafsson/jgosonnet/internal/ast"
)

type TraceSignal struct {
	Str  string
	Rest Value
}

func (t *TraceSignal) Error() string {
	return "TRACE: " + t.Str
}

func lookupProgram(bp unsafe.Pointer, ctx Context) *fast.Program {
	if ctx.State == nil || ctx.State.Environment == nil || ctx.State.Environment.Importer == nil {
		return nil
	}
	return ctx.State.Environment.Importer.LookupProgram(bp)
}

func (t *TraceSignal) Print(w io.Writer, bp unsafe.Pointer, offset uint32, ctx Context) error {
	prog := lookupProgram(bp, ctx)
	if prog == nil {
		_, err := fmt.Fprintf(w, "TRACE: %s\n", t.Str)
		return err
	}
	loc, _ := prog.NodeInfo(offset)
	if !loc.IsSet() {
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

func WrapError(err error, bp unsafe.Pointer, offset uint32, ctx Context) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*TraceSignal); ok {
		return err
	}

	traceErr, ok := err.(*TraceError)
	if !ok {
		traceErr = &TraceError{
			Err:    err,
			Frames: make([]Frame, 0, 8),
		}
	}

	prog := lookupProgram(bp, ctx)
	if prog == nil {
		return traceErr
	}

	loc, ctxName := prog.NodeInfo(offset)
	if !loc.IsSet() {
		return traceErr
	}

	pos := loc.String()
	if pos == "" {
		return traceErr
	}

	frame := Frame{Pos: pos, Name: ctxName}

	// always append first frame or when new context boundary.
	numFrames := len(traceErr.Frames)
	if numFrames == 0 || traceErr.Frames[numFrames-1].Name != frame.Name {
		traceErr.Frames = append(traceErr.Frames, frame)
	}

	return traceErr
}
