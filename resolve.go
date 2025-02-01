package flagset

import (
	"errors"
	"flag"
	"strconv"
	"strings"
	"time"

	er "github.com/daved/flagset/fserrs"
	"github.com/daved/flagset/vtype"
)

func resolve(flags []*Flag, args []string) ([]string, error) {
	var flag *Flag

	for i, arg := range args {
		if flag != nil { // expecting flag value
			if err := hydrate(flag.val, arg); err != nil {
				return nil, err
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
				flag, err = findFlag(flags, name)
				if err != nil {
					return nil, err
				}

				isBool, boolErr := resolveBool(flag)
				if boolErr != nil {
					return nil, boolErr
				}
				if isBool {
					flag = nil
				}
				continue
			}

			var raw string
			name, raw, _ = strings.Cut(name, "=")
			flag, err = findFlag(flags, name)
			if err != nil {
				return nil, err
			}

			if err := hydrate(flag.val, raw); err != nil {
				return nil, err
			}

			flag = nil

		default: // short flag
			name := arg[1:2]

			var err error
			flag, err = findFlag(flags, name)
			if err != nil {
				return nil, err
			}

			isBool, boolErr := resolveBool(flag)
			if boolErr != nil {
				return nil, boolErr
			}
			if isBool {
				flag = nil
			}
		}
	}

	return nil, nil
}

func findFlag(flags []*Flag, name string) (*Flag, error) {
	for _, flag := range flags {
		ss := flag.shorts
		if len(name) > 1 {
			ss = flag.longs
		}
		if sliceContains(ss, name) {
			return flag, nil
		}
	}
	return nil, errors.New("unknown flag name")
}

func sliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func resolveBool(f *Flag) (bool, error) {
	switch v := f.val.(type) {
	case *bool:
		*v = true
		return true, nil

	case vtype.FlagCallback:
		if v.IsBool() {
			err := v.OnFlag("")
			return true, err
		}
	}

	return false, nil
}

func hydrate(val any, raw string) error {
	newError := er.NewParseError

	switch v := val.(type) {
	case *string:
		*v = raw

	case *bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return newError(er.NewConvertRawError(err))
		}
		*v = b

	case *int:
		n, err := strconv.Atoi(raw)
		if err != nil {
			return newError(er.NewConvertRawError(err))
		}
		*v = n

	case *int64:
		n, err := strconv.ParseInt(raw, 10, 0)
		if err != nil {
			return newError(er.NewConvertRawError(err))
		}
		*v = n

	case *uint:
		n, err := strconv.ParseUint(raw, 10, 0)
		if err != nil {
			return newError(er.NewConvertRawError(err))
		}
		*v = uint(n)

	case *uint64:
		n, err := strconv.ParseUint(raw, 10, 0)
		if err != nil {
			return newError(er.NewConvertRawError(err))
		}
		*v = n

	case *float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return newError(er.NewConvertRawError(err))
		}
		*v = f

	case *time.Duration:
		d, err := time.ParseDuration(raw)
		if err != nil {
			return newError(er.NewConvertRawError(err))
		}
		*v = d

	case vtype.TextMarshalUnmarshaler:
		if err := v.UnmarshalText([]byte(raw)); err != nil {
			return newError(er.NewConvertRawError(err))
		}

	case flag.Value:
		if err := v.Set(raw); err != nil {
			return newError(er.NewConvertRawError(err))
		}

	case vtype.FlagCallback:
		if err := v.OnFlag(raw); err != nil {
			return newError(er.NewConvertRawError(err))
		}
	}

	return nil
}
