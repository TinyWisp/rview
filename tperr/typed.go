package tperr

import (
	"fmt"
)

type Typed interface {
	Is(etype string)
	GetEtype() string
	GetVars() []any
}

type TypedError struct {
	etype string
	vars  []any
}

func NewTypedError(etype string, vars ...any) *TypedError {
	return &TypedError{
		etype: etype,
		vars:  vars,
	}
}

func (te *TypedError) Error() string {
	return fmt.Sprintf(T(te.etype), te.vars...)
}

func (te *TypedError) Is(etype string) bool {
	return te.etype == etype
}

func (te *TypedError) GetEtype() string {
	return te.etype
}

func (te *TypedError) GetVars() []any {
	return te.vars
}
