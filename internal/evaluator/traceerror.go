package evaluator

import (
	"fmt"
	"strings"

	"github.com/google/go-jsonnet/ast"
)

type TraceError struct {
	Err    error
	Frames []Frame
}

func (t *TraceError) Error() string {
	var b strings.Builder

	b.WriteString(t.Err.Error())
	b.WriteByte('\n')

	for _, frame := range t.Frames {
		fmt.Fprintf(&b, "\t%s\t%s\n", frame.Pos, frame.Name)
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

	frame := Frame{Pos: node.Loc().String()}
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
