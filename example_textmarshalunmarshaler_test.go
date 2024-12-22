package flagset_test

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/daved/flagset"
)

type Color int

const (
	Red Color = iota + 1
	Green
	Blue
)

func (c Color) MarshalText() (text []byte, err error) {
	switch c {
	case Red:
		return []byte("red"), nil
	case Green:
		return []byte("green"), nil
	case Blue:
		return []byte("blue"), nil
	default:
		return nil, errors.New("invalid color: " + strconv.Itoa(int(c)))
	}
}

func (c *Color) UnmarshalText(text []byte) error {
	switch s := string(text); s {
	case "red":
		*c = Red
	case "green":
		*c = Green
	case "blue":
		*c = Blue
	default:
		return errors.New("invalid color: " + s)
	}
	return nil
}

func ExampleTextMarshalUnmarshaler() {
	c := Red

	fs := flagset.New("app")
	fs.Flag(&c, "color|c", "Color to use.")

	args := []string{"--color=green"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		return
	}

	colorTxt, _ := c.MarshalText()
	fmt.Printf("Color: %s\n", colorTxt)
	fmt.Println()
	fmt.Println(fs.Usage())

	// Output:
	// Color: green
	//
	// Flags for app:
	//
	//     -c, --color  =VALUE    default: red
	//         Color to use.
}
