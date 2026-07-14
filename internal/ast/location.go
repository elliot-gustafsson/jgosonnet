package ast

import "fmt"

type NodeContext struct {
	Context *string
	File    string
	Begin   Location
	End     Location
}

type Location struct {
	Line   uint32
	Column uint32
}

func (t NodeContext) Empty() bool {
	return t.Begin.Line == 0
}

func (t NodeContext) String() string {

	var filePrefix string
	if len(t.File) > 0 {
		filePrefix = t.File + ":"
	}

	if t.Begin.Line == t.End.Line {
		if t.Begin.Column == t.End.Column {
			return fmt.Sprintf("%s%v", filePrefix, t.Begin.String())
		}
		return fmt.Sprintf("%s%v-%v", filePrefix, t.Begin.String(), t.End.Column)
	}

	return fmt.Sprintf("%s(%v)-(%v)", filePrefix, t.Begin.String(), t.End.String())
}

func (l Location) String() string {
	return fmt.Sprintf("%v:%v", l.Line, l.Column)
}
