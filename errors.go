package flagset

import (
	"github.com/daved/flagset/fserrs"
)

// Error types forward basic error types from the fserrs package for access and
// documentation. If an error has interesting behavior, it should be defined
// directly in this package.
type (
	HydrateError          = fserrs.HydrateError
	UnrecognizedFlagError = fserrs.UnrecognizedFlagError
)
