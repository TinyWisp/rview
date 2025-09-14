package tperr

import (
	"fmt"
	"strings"
)

type TraceableTypedError struct {
	TypedError
	pos int
	raw string
}

func (tte *TraceableTypedError) Error() string {
	msg := ""

	msg = fmt.Sprintf(T(tte.etype), tte.vars...) + "\n"

	if tte.pos == -1 || tte.raw == "" {
		tte.pos = len(tte.raw) - 1
	}

	lines := strings.Split(tte.raw, "\n")
	lastLineEndPos := 0
	errRow := 0
	errCol := tte.pos
	for idx := 0; idx < len(lines); idx++ {
		if lastLineEndPos+len(lines[idx]) >= tte.pos {
			errRow = idx
			errCol = tte.pos - lastLineEndPos
			break
		}
		lastLineEndPos += len(lines[idx]) + 1
	}

	beginLine := errRow - 2
	if beginLine < 0 {
		beginLine = 0
	}
	for idx := beginLine; idx <= errRow; idx++ {
		msg += lines[idx] + "\n"
	}
	for idx := 0; idx < errCol; idx++ {
		msg += " "
	}
	msg += "^\n"
	for idx := errRow + 1; idx < errRow+2 && idx < len(lines); idx++ {
		msg += lines[idx] + "\n"
	}

	// fmt.Printf("errRow:%d, errCol:%d, beginLine:%d, endLine:%d\n", errRow, errCol, beginLine, len(lines)-1)
	// fmt.Println(lines)

	return msg
}

func (tte *TraceableTypedError) AddOffset(offset int) {
	tte.pos += offset
}

func (tte *TraceableTypedError) SetRaw(raw string) {
	tte.raw = raw
}

func (tte *TraceableTypedError) SetPos(pos int) {
	tte.pos = pos
}

func (tte *TraceableTypedError) IsKindOf(kind string) bool {
	return strings.HasPrefix(tte.etype, kind)
}

func NewTraceableTypedError(raw string, pos int, etype string, vars ...any) *TraceableTypedError {
	return &TraceableTypedError{
		TypedError: TypedError{
			etype: etype,
			vars:  vars,
		},
		raw: raw,
		pos: pos,
	}
}

func IsErrorType(err error, etype string) bool {
	if terr, ok := err.(*TypedError); ok {
		return terr.Is(etype)
	}

	if terr, ok := err.(*TraceableTypedError); ok {
		return terr.Is(etype)
	}

	return false
}
