package flagset

import (
	"fmt"
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
				flag, err = findFlag(flags, name)
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
			flag, err = findFlag(flags, name)
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
			flag, err = findFlag(flags, name)
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
		if err := vtype.Hydrate(flag, ""); err != nil {
			return nil, wrap(err, flag.name)
		}
	}

	return nil, nil
}

func findFlag(flags []*Flag, name string) (*namedFlag, error) {
	for _, flag := range flags {
		ss := flag.shorts
		if len(name) > 1 {
			ss = flag.longs
		}
		if sliceContains(ss, name) {
			return &namedFlag{flag, name}, nil
		}
	}
	return nil, fmt.Errorf("find flag: %w", er.ErrUnrecognizedFlag)
}

func sliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
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
