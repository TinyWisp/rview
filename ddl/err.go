package ddl

import (
	"fmt"
	"strings"

	"github.com/TinyWisp/rview/tperr"
)

type DdlError struct {
	pos   int
	ddl   string
	etype string
	vars  []any
}

func NewDdlError(ddl string, pos int, etype string, vars ...any) *DdlError {
	return &DdlError{
		ddl:   ddl,
		etype: etype,
		pos:   pos,
		vars:  vars,
	}
}

func (de *DdlError) Error() string {
	msg := ""

	msg = fmt.Sprintf(tperr.T(de.etype), de.vars...) + "\n"

	if de.pos < 0 || de.ddl == "" {
		return msg
	}

	lines := strings.Split(de.ddl, "\n")
	lastLineEndPos := 0
	errRow := 0
	errCol := de.pos
	for idx := 0; idx < len(lines); idx++ {
		if lastLineEndPos+len(lines[idx]) >= de.pos {
			errRow = idx
			errCol = de.pos - lastLineEndPos
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

	//fmt.Printf("errRow:%d, errCol:%d, beginLine:%d, endLine:%d\n", errRow, errCol, beginLine, len(lines)-1)
	//fmt.Println(lines)

	return msg
}

func (de *DdlError) AddOffset(offset int) {
	de.pos += offset
}

func (de *DdlError) SetDdl(ddl string) {
	de.ddl = ddl
}

func (de *DdlError) SetPos(pos int) {
	de.pos = pos
}

func (de *DdlError) Is(etype string) bool {
	return etype == de.etype
}

func (de *DdlError) IsExpError() bool {
	return strings.HasPrefix(de.etype, "exp")
}

func (de *DdlError) IsCssError() bool {
	return strings.HasPrefix(de.etype, "css")
}
