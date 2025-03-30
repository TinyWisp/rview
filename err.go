package rview

import (
	"fmt"

	"github.com/TinyWisp/rview/i18n"
)

type FmtError struct {
	err  string
	vars []any
}

func NewError(err string, vars ...any) FmtError {
	return FmtError{
		err:  err,
		vars: vars,
	}
}

func (fe FmtError) Error() string {
	return fmt.Sprintf(i18n.T(fe.err), fe.vars...)
}
