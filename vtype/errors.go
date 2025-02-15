package vtype

import (
	"errors"
	"fmt"
	"reflect"
)

type Error struct {
	child error
}

func NewError(child error) *Error {
	return &Error{child}
}

func (e *Error) Error() string {
	return fmt.Sprintf("vtype: %v", e.child)
}

func (e *Error) Unwrap() error {
	return e.child
}

func (e *Error) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type HydrateError struct {
	child error
	Val   any
}

func NewHydrateError(child error, val any) *HydrateError {
	return &HydrateError{child, val}
}

func (e *HydrateError) Error() string {
	return fmt.Sprintf("hydrate (type: %T): %v", e.Val, e.child)
}

func (e *HydrateError) Unwrap() error {
	return e.child
}

func (e *HydrateError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

var ErrUnsupportedType = errors.New("unsupported type")
