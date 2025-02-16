package fserrs

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
	return fmt.Sprintf("flagset: %v", e.child)
}

func (e *Error) Unwrap() error {
	return e.child
}

func (e *Error) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

var ErrFlagUnrecognized = errors.New("flag unrecognized")

type ParseError struct {
	child error
}

func NewParseError(child error) *ParseError {
	return &ParseError{child}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse: %v", e.child)
}

func (e *ParseError) Unwrap() error {
	return e.child
}

func (e *ParseError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type ResolveError struct {
	child    error
	FlagName string
}

func NewResolveError(child error, flagName string) *ResolveError {
	return &ResolveError{child, flagName}
}

func (e *ResolveError) Error() string {
	return fmt.Sprintf("resolve (flag name: %s): %v", e.FlagName, e.child)
}

func (e *ResolveError) Unwrap() error {
	return e.child
}

func (e *ResolveError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}
