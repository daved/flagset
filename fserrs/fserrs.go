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
