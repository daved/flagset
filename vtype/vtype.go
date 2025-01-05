package vtype

import (
	"encoding"
	"strconv"
)

// TextMarshalUnmarshaler descibes types that are able to be marshaled to and
// unmarshaled from text.
type TextMarshalUnmarshaler interface {
	encoding.TextUnmarshaler
	encoding.TextMarshaler
}

type FlagCallback interface {
	OnFlag(val string) error
	IsBool() bool
}

// FlagFunc describes functions that can be called when a flag option is
// succesfully parsed. Currently, this cannot pass errors values back to callers
// as the stdlib flag pkg eats them.
type FlagFunc func(string) error

func (f FlagFunc) OnFlag(val string) error {
	return f(val)
}

func (f FlagFunc) IsBool() bool { return false }

// FlagBoolFunc describes functions that can be called when a bool flag option
// is succesfully parsed. Currently, this cannot pass errors values back to
// callers as the stdlib flag pkg eats them.
type FlagBoolFunc func(bool) error

func (f FlagBoolFunc) OnFlag(s string) error {
	if s == "" {
		return f(true)
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	return f(b)
}

func (f FlagBoolFunc) IsBool() bool { return true }
