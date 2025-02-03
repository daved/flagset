package flagset

import (
	"fmt"
	"reflect"

	"github.com/daved/flagset/fserrs"
)

// Error types forward basic error types from the fserrs package for access and
// documentation. If an error has interesting behavior, it should be defined
// directly in this package.
type (
	ParseError = fserrs.ParseError
)

type HydrateError struct {
	Flag  string
	child error
}

func NewHydrateError(flag string, child error) *HydrateError {
	return &HydrateError{flag, child}
}

func (e *HydrateError) Error() string {
	return fmt.Sprintf("hydrate (%s): %v", e.Flag, e.child)
}

func (e *HydrateError) Unwrap() error {
	return e.child
}

func (e *HydrateError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type UnrecognizedFlagError struct {
	Arg string
}

func NewUnrecognizedFlagError(arg string) *UnrecognizedFlagError {
	return &UnrecognizedFlagError{arg}
}

func (e *UnrecognizedFlagError) Error() string {
	return fmt.Sprintf("unrecognized flag: %q", e.Arg)
}

func (e *UnrecognizedFlagError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}
