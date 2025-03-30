package ddl

import (
	"fmt"
	"strings"

	"github.com/TinyWisp/rview/i18n"
)

type DdlParseError struct {
	pos int
	ddl string
	err string
}

func NewDdlParseError(ddl string, err string, pos int) *DdlParseError {
	return &DdlParseError{
		ddl: ddl,
		err: err,
		pos: pos,
	}
}

func (dpe *DdlParseError) Error() string {
	msg := ""

	msg = i18n.T(dpe.err) + "\n"

	if dpe.pos == -1 {
		dpe.pos = len(dpe.ddl) - 1
	}

	lines := strings.Split(dpe.ddl, "\n")
	lastLineEndPos := 0
	errRow := 0
	errCol := dpe.pos
	for idx := 0; idx < len(lines); idx++ {
		if lastLineEndPos+len(lines[idx]) >= dpe.pos {
			errRow = idx
			errCol = dpe.pos - lastLineEndPos
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

	fmt.Printf("errRow:%d, errCol:%d, beginLine:%d, endLine:%d\n", errRow, errCol, beginLine, len(lines)-1)
	fmt.Println(lines)

	return msg
}

func (dpe *DdlParseError) AddOffset(offset int) {
	dpe.pos += offset
}

func (dpe *DdlParseError) SetDdl(ddl string) {
	dpe.ddl = ddl
}

func (dpe *DdlParseError) SetPos(pos int) {
	dpe.pos = pos
}

func (dpe *DdlParseError) IsExpError() bool {
	return strings.HasPrefix(dpe.err, "exp")
}

func (dpe *DdlParseError) IsCssError() bool {
	return strings.HasPrefix(dpe.err, "css")
}
