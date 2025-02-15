package flagset

import (
	"github.com/daved/flagset/fserrs"
	"github.com/daved/flagset/vtype"
)

// Error types forward basic error types from sub-packages for access and
// documentation. If an error has interesting behavior, it should be defined
// directly in this package.
type (
	Error        = fserrs.Error
	ParseError   = fserrs.ParseError
	ResolveError = fserrs.ResolveError
)

var (
	ErrUnrecognizedFlag = fserrs.ErrUnrecognizedFlag
	ErrUnsupportedType  = vtype.ErrUnsupportedType
)
