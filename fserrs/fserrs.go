package fserrs

import (
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
	child error
}

func NewResolveError(child error) *ResolveError {
	return &ResolveError{child}
}

func (e *ResolveError) Error() string {
	return fmt.Sprintf("resolve: %v", e.child)
}

func (e *ResolveError) Unwrap() error {
	return e.child
}

func (e *ResolveError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type FindFlagError struct {
	child error
}

func NewFindFlagError(child error) *FindFlagError {
	return &FindFlagError{child}
}

func (e *FindFlagError) Error() string {
	return fmt.Sprintf("find flag: %v", e.child)
}

func (e *FindFlagError) Unwrap() error {
	return e.child
}

func (e *FindFlagError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type FlagHydrateError struct {
	Name  string
	child error
}

func NewHydrateError(name string, child error) *FlagHydrateError {
	return &FlagHydrateError{name, child}
}

func (e *FlagHydrateError) Error() string {
	return fmt.Sprintf("hydrate (for %q): %v", e.Name, e.child)
}

func (e *FlagHydrateError) Unwrap() error {
	return e.child
}

func (e *FlagHydrateError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type FlagUnrecognizedError struct {
	Name string
}

func NewFlagUnrecognizedError(name string) *FlagUnrecognizedError {
	return &FlagUnrecognizedError{name}
}

func (e *FlagUnrecognizedError) Error() string {
	return fmt.Sprintf("unrecognized flag: %q", e.Name)
}

func (e *FlagUnrecognizedError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}
