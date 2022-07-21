package echotool

import (
	"fmt"
	"sync"
)

type Kind int

const (
	KindHeader Kind = iota
	KindParam
	KindQuery
	KindPostForm
	KindForm
)

var _ fmt.Stringer = (*Kind)(nil)

func (k Kind) String() string {
	switch k {
	case KindHeader:
		return "header"
	case KindParam:
		return "param"
	case KindQuery:
		return "query"
	case KindPostForm:
		return "postform"
	case KindForm:
		return "form"
	}
	return fmt.Sprintf("unknown kind: %d", k)
}

type EmptyError struct {
	kind Kind
	key  string
}

var _ error = (*EmptyError)(nil)

func NewEmptyError(kind Kind, key string) *EmptyError {
	return &EmptyError{
		kind: kind,
		key:  key,
	}
}

func NewHeaderEmptyError(key string) *EmptyError {
	return NewEmptyError(KindHeader, key)
}

func NewParamEmptyError(key string) *EmptyError {
	return NewEmptyError(KindParam, key)
}

func NewQueryEmptyError(key string) *EmptyError {
	return NewEmptyError(KindQuery, key)
}

func NewPostFormEmptyError(key string) *EmptyError {
	return NewEmptyError(KindPostForm, key)
}

func NewFormEmptyError(key string) *EmptyError {
	return NewEmptyError(KindForm, key)
}

func (e EmptyError) Error() string {
	return fmt.Sprintf("%s(%s) is empty", e.key, e.kind.String())
}

func IsEmptyError(err error) bool {
	_, ok := err.(*EmptyError)
	return ok
}

type EchotoolError struct {
	code int
	err  error
}

var _ error = (*EchotoolError)(nil)

func (e EchotoolError) GetCode() int {
	return e.code
}

func (e EchotoolError) GetError() error {
	return e.err
}

func (e EchotoolError) Error() string {
	return fmt.Sprintf("%s - %+v", CodeMsg(e.code), e.err)
}

func (e *EchotoolError) reset() {
	e.code = 0
	e.err = nil
}

var echotoolErrorPool sync.Pool

func AcquireEchotoolError(code int, err error) (e *EchotoolError) {
	if v := echotoolErrorPool.Get(); v != nil {
		e = v.(*EchotoolError)
	} else {
		e = &EchotoolError{}
	}

	e.code = code
	e.err = err
	return
}

func ReleaseEchotoolError(e *EchotoolError) {
	if e != nil {
		e.reset()
		echotoolErrorPool.Put(e)
	}
}

func IsEchotoolError(err error) bool {
	_, ok := err.(*EchotoolError)
	return ok
}
