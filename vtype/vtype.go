package vtype

import (
	"bytes"
	"encoding"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// TextMarshalUnmarshaler descibes types that are able to be marshaled to and
// unmarshaled from text.
type TextMarshalUnmarshaler interface {
	encoding.TextUnmarshaler
	encoding.TextMarshaler
}

// FlagCallback describes types that will have a callback ("OnFlag") run when
// the associated flag is parsed. The IsBool method is necessary because bool
// flags are treated specially during parsing. While it is reasonable to
// implement this type directly, it is recommended to, instead, provide a
// function compatible with [FlagFunc] or [FlagBoolFunc] as conversion is done
// automatically, and both Func types are implementations of FlagCallback.
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

type Slice struct {
	ptrToSlice any
	started    bool

	TypeName  string
	SplitEach bool
	Separator string
}

func MakeSlice(ptrToSlice any) Slice {
	return Slice{
		ptrToSlice: ptrToSlice,
		Separator:  ",",
	}
}

func (s *Slice) UnmarshalText(text []byte) error {
	vo := reflect.ValueOf(s.ptrToSlice)
	isPtr := vo.Kind() == reflect.Pointer
	if isPtr {
		vo = vo.Elem()
	}
	if !isPtr || vo.Kind() != reflect.Slice {
		return errors.New("slice: contained value is not a pointer to a slice")
	}

	if !s.started {
		slice := reflect.MakeSlice(vo.Type(), 0, 0)
		reflect.ValueOf(s.ptrToSlice).Elem().Set(slice)
	}
	s.started = true

	valType := vo.Type().Elem()

	sep := s.Separator
	if !s.SplitEach {
		sep = "<><>"
	}

	for _, chunk := range bytes.Split(text, []byte(sep)) {
		item := reflect.New(valType)
		if err := Hydrate(item.Interface(), string(chunk)); err != nil {
			return fmt.Errorf("slice: unmarshal text: %w", err)
		}

		slice := reflect.Append(vo, item.Elem())
		reflect.ValueOf(s.ptrToSlice).Elem().Set(slice)
	}

	return nil
}

func (s *Slice) MarshalText() ([]byte, error) {
	vo := reflect.ValueOf(s.ptrToSlice)
	isPtr := vo.Kind() == reflect.Pointer
	if isPtr {
		vo = vo.Elem()
	}
	if !isPtr || vo.Kind() != reflect.Slice {
		return nil, errors.New("slice: contained value is not a pointer to a slice")
	}

	out := make([]string, vo.Len())
	for i := 0; i < vo.Len(); i++ {
		out[i] = fmt.Sprint(vo.Index(i).Interface())
	}
	return []byte(strings.Join(out, s.Separator)), nil
}

func (s *Slice) UsageTypeName() string {
	rv := reflect.ValueOf(s.ptrToSlice)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	name := rv.Type().Elem().Name()

	if s.SplitEach {
		name += "(csv)"
	}

	return name
}

func (s *Slice) IsBool() bool {
	return reflect.ValueOf(s.ptrToSlice).Elem().Type().Elem().Kind() == reflect.Bool
}

func Hydrate(val any, raw string) error {
	wrap := func(err error) error {
		return NewError(NewHydrateError(err, val))
	}

	switch v := val.(type) {
	case error:
		return wrap(v)

	case *string:
		*v = raw

	case *bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return wrap(err)
		}
		*v = b

	case *int:
		n, err := strconv.Atoi(raw)
		if err != nil {
			return wrap(err)
		}
		*v = n

	case *int64:
		n, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return wrap(err)
		}
		*v = n

	case *int8:
		n, err := strconv.ParseInt(raw, 10, 8)
		if err != nil {
			return wrap(err)
		}
		*v = int8(n)

	case *int16:
		n, err := strconv.ParseInt(raw, 10, 16)
		if err != nil {
			return wrap(err)
		}
		*v = int16(n)

	case *int32:
		n, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return wrap(err)
		}
		*v = int32(n)

	case *uint:
		n, err := strconv.ParseUint(raw, 10, 0)
		if err != nil {
			return wrap(err)
		}
		*v = uint(n)

	case *uint64:
		n, err := strconv.ParseUint(raw, 10, 0)
		if err != nil {
			return wrap(err)
		}
		*v = n

	case *uint8:
		n, err := strconv.ParseUint(raw, 10, 8)
		if err != nil {
			return wrap(err)
		}
		*v = uint8(n)

	case *uint16:
		n, err := strconv.ParseUint(raw, 10, 16)
		if err != nil {
			return wrap(err)
		}
		*v = uint16(n)

	case *uint32:
		n, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return wrap(err)
		}
		*v = uint32(n)

	case *float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return wrap(err)
		}
		*v = f

	case *float32:
		f, err := strconv.ParseFloat(raw, 32)
		if err != nil {
			return wrap(err)
		}
		*v = float32(f)

	case *time.Duration:
		d, err := time.ParseDuration(raw)
		if err != nil {
			return wrap(err)
		}
		*v = d

	case TextMarshalUnmarshaler:
		if err := v.UnmarshalText([]byte(raw)); err != nil {
			return wrap(err)
		}

	case flag.Value:
		if err := v.Set(raw); err != nil {
			return wrap(err)
		}

	case FlagCallback:
		if err := v.OnFlag(raw); err != nil {
			return wrap(err)
		}

	default:
		return wrap(ErrUnsupportedType)
	}

	return nil
}
