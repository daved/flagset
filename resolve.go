package flagset

import (
	"flag"
	"strconv"
	"strings"
	"time"

	er "github.com/daved/flagset/fserrs"
	"github.com/daved/flagset/vtype"
)

type namedFlag struct {
	flag *Flag
	name string
}

func resolve(flags []*Flag, args []string) ([]string, error) {
	wrap := er.NewResolveError
	var flag *namedFlag

	for i, arg := range args {
		if flag != nil { // expecting flag value
			if err := hydrate(flag, arg); err != nil {
				return nil, wrap(err)
			}

			flag = nil
			continue
		}

		// expecting flag or operand
		switch {
		case arg[0] != '-' || arg == "-": // operand
			return args[i:], nil

		case arg == "--": // end resolution
			if i+1 < len(args) {
				return args[i+1:], nil
			}
			return nil, nil

		case arg[0:2] == "--": // long flag
			var err error

			name := arg[2:]
			if !strings.Contains(name, "=") {
				flag, err = findNamedFlag(flags, name)
				if err != nil {
					return nil, wrap(err)
				}

				isBool, boolErr := hydrateBool(flag)
				if boolErr != nil {
					return nil, wrap(boolErr)
				}
				if isBool {
					flag = nil
				}
				continue
			}

			var raw string
			name, raw, _ = strings.Cut(name, "=")
			flag, err = findNamedFlag(flags, name)
			if err != nil {
				return nil, wrap(err)
			}

			if err := hydrate(flag, raw); err != nil {
				return nil, wrap(err)
			}

			flag = nil

		default: // short flag
			name := arg[1:2]

			var err error
			flag, err = findNamedFlag(flags, name)
			if err != nil {
				return nil, wrap(err)
			}

			isBool, boolErr := hydrateBool(flag)
			if boolErr != nil {
				return nil, wrap(boolErr)
			}
			if isBool {
				flag = nil
			}
		}
	}

	if flag != nil {
		if err := hydrate(flag, ""); err != nil {
			return nil, wrap(err)
		}
	}

	return nil, nil
}

func findNamedFlag(flags []*Flag, name string) (*namedFlag, error) {
	for _, flag := range flags {
		ss := flag.shorts
		if len(name) > 1 {
			ss = flag.longs
		}
		if sliceContains(ss, name) {
			return &namedFlag{flag, name}, nil
		}
	}
	return nil, er.NewFindFlagError(er.NewFlagUnrecognizedError(name))
}

func sliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func hydrateBool(fa *namedFlag) (isBool bool, boolErr error) {
	wrap := er.NewHydrateError

	switch v := fa.flag.val.(type) {
	case *bool:
		*v = true
		return true, nil

	case vtype.FlagCallback:
		if v.IsBool() {
			// when isBool, OnFlag won't be called elsewhere
			err := v.OnFlag("")
			return true, wrap(fa.name, err)
		}
		return false, nil

	case error: // always treat as bool
		return true, wrap(fa.name, v)

	default:
		return false, nil
	}
}

func hydrate(fa *namedFlag, raw string) error {
	wrap := er.NewHydrateError

	switch v := fa.flag.val.(type) {
	case error:
		return wrap(fa.name, v)

	case *string:
		*v = raw

	case *bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = b

	case *int:
		n, err := strconv.Atoi(raw)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = n

	case *int64:
		n, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = n

	case *int8:
		n, err := strconv.ParseInt(raw, 10, 8)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = int8(n)

	case *int16:
		n, err := strconv.ParseInt(raw, 10, 16)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = int16(n)

	case *int32:
		n, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = int32(n)

	case *uint:
		n, err := strconv.ParseUint(raw, 10, 0)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = uint(n)

	case *uint64:
		n, err := strconv.ParseUint(raw, 10, 0)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = n

	case *uint8:
		n, err := strconv.ParseUint(raw, 10, 8)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = uint8(n)

	case *uint16:
		n, err := strconv.ParseUint(raw, 10, 16)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = uint16(n)

	case *uint32:
		n, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = uint32(n)

	case *float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = f

	case *float32:
		f, err := strconv.ParseFloat(raw, 32)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = float32(f)

	case *time.Duration:
		d, err := time.ParseDuration(raw)
		if err != nil {
			return wrap(fa.name, err)
		}
		*v = d

	case vtype.TextMarshalUnmarshaler:
		if err := v.UnmarshalText([]byte(raw)); err != nil {
			return wrap(fa.name, err)
		}

	case flag.Value:
		if err := v.Set(raw); err != nil {
			return wrap(fa.name, err)
		}

	case vtype.FlagCallback:
		if err := v.OnFlag(raw); err != nil {
			return wrap(fa.name, err)
		}
	}

	return nil
}
