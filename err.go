package rview

import (
	"fmt"

	"github.com/TinyWisp/rview/i18n"
)

type FmtError struct {
	etype string
	vars  []any
}

func NewError(etype string, vars ...any) *FmtError {
	return &FmtError{
		etype: etype,
		vars:  vars,
	}
}

func (fe *FmtError) Error() string {
	return fmt.Sprintf(i18n.T(fe.etype), fe.vars...)
}

func (fe *FmtError) Is(etype string) bool {
	return fe.etype == etype
}

func IsErrorType(err error, etype string) bool {
	if ferr, ok := err.(*FmtError); ok {
		return ferr.Is(etype)
	}

	return false
}
