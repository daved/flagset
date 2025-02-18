package flagset

import (
	"reflect"
	"strings"

	er "github.com/daved/flagset/fserrs"
	"github.com/daved/vtype"
)

type namedFlag struct {
	flag *Flag
	name string
}

func resolveFlags(flags []*Flag, args []string) ([]string, error) {
	wrap := er.NewResolveError
	var flag *namedFlag

	for i, arg := range args {
		if flag != nil { // expecting flag value
			if err := vtype.Hydrate(flag.flag.val, arg); err != nil {
				return nil, wrap(err, flag.name)
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
				flag, err = lookupFlagAsNamedFlag(flags, name)
				if err != nil {
					return nil, wrap(err, name)
				}

				if bv, ok := boolValRaw(flag.flag.val); ok {
					if err = vtype.Hydrate(flag.flag.val, bv); err != nil {
						return nil, wrap(err, flag.name)
					}
					flag = nil
				}
				continue
			}

			var raw string
			name, raw, _ = strings.Cut(name, "=")
			flag, err = lookupFlagAsNamedFlag(flags, name)
			if err != nil {
				return nil, wrap(err, name)
			}

			if err := vtype.Hydrate(flag.flag.val, raw); err != nil {
				return nil, wrap(err, flag.name)
			}

			flag = nil

		default: // short flag
			name := arg[1:2]

			var err error
			flag, err = lookupFlagAsNamedFlag(flags, name)
			if err != nil {
				return nil, wrap(err, name)
			}

			if bv, ok := boolValRaw(flag.flag.val); ok {
				if err = vtype.Hydrate(flag.flag.val, bv); err != nil {
					return nil, wrap(err, flag.name)
				}
				flag = nil
			}
		}
	}

	if flag != nil {
		if err := vtype.Hydrate(flag.flag.val, ""); err != nil {
			return nil, wrap(err, flag.name)
		}
	}

	return nil, nil
}

func lookupFlagAsNamedFlag(flags []*Flag, name string) (*namedFlag, error) {
	if flag := lookupFlag(flags, name); flag != nil {
		return &namedFlag{
			flag: flag,
			name: name,
		}, nil
	}
	return nil, er.ErrFlagUnrecognized
}

func boolValRaw(val any) (string, bool) {
	switch v := val.(type) {
	case interface{ IsBool() bool }:
		if reflect.ValueOf(val).Kind() == reflect.Func {
			return "", v.IsBool()
		}
		return "true", v.IsBool()

	case *bool:
		return "true", true

	case error:
		return "", true

	default:
		return "", false
	}
}
