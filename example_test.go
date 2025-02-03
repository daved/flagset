package flagset_test

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/daved/flagset"
)

func Example() {
	var (
		info    = "default-value"
		num     int
		verbose bool
	)

	fs := flagset.New("app")
	fs.Flag(&info, "info|i", "Interesting info.")
	fs.Flag(&num, "num|n", "Number with no usage.").HideUsage = true
	fs.Flag(&verbose, "verbose|v", "Set verbose output.")

	args := []string{"--info=non-default", "-n", "42", "-v"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fs.Usage())
	fmt.Printf("Info: %s, Num: %d, Verbose: %t\n", info, num, verbose)
	// Output:
	// Flags for app:
	//
	//     -i, --info  =STRING    default: default-value
	//         Interesting info.
	//
	//     -v, --verbose  [=BOOL]    default: false
	//         Set verbose output.
	//
	// Info: non-default, Num: 42, Verbose: true
}

func Example_flagFunc() {
	onDo := func(flagValue string) error {
		fmt.Printf("Func called: Doing %s\n\n", flagValue)
		return nil
	}

	fs := flagset.New("app")
	fs.Flag(onDo, "do|d", "Run callback.")

	args := []string{"--do=something"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fs.Usage())
	// Output:
	// Func called: Doing something
	//
	// Flags for app:
	//
	//     -d, --do  =VALUE
	//         Run callback.
}

type URLValue struct {
	URL *url.URL
}

func (v URLValue) String() string {
	if v.URL != nil {
		return v.URL.String()
	}
	return ""
}

func (v URLValue) Set(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	*v.URL = *u
	return nil
}

func Example_flagValue() {
	url := URLValue{&url.URL{}}

	fs := flagset.New("app")
	fs.Flag(url, "url|u", "URL to use.")

	args := []string{"--url=https://example.com"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fs.Usage())
	fmt.Printf("URL: %s\n", url)
	// Output:
	// Flags for app:
	//
	//     -u, --url  =VALUE
	//         URL to use.
	//
	// URL: https://example.com
}

func Example_helpError() {
	var (
		debug   bool
		errHelp = errors.New("help requested")
		verbose bool
	)

	fs := flagset.New("app")
	fs.Flag(&debug, "debug|d", "Set debug output.")
	fs.Flag(errHelp, "help|h", "Display usage output.")
	fs.Flag(&verbose, "verbose|v", "Set verbose output.")

	args := []string{"-d", "-h", "-v"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(fs.Usage())

		if !errors.Is(err, errHelp) {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// help is an error value, so parsing will halt at -h
	fmt.Printf("Debug: %t, Verbose: %t\n", debug, verbose)
	// Output:
	// Flags for app:
	//
	//     -d, --debug  [=BOOL]    default: false
	//         Set debug output.
	//
	//     -h, --help
	//         Display usage output.
	//
	//     -v, --verbose  [=BOOL]    default: false
	//         Set verbose output.
	//
	// Debug: true, Verbose: false
}

type Color int

const (
	Unset Color = iota
	Green
	Blue
)

func (c Color) MarshalText() (text []byte, err error) {
	switch c {
	case Unset:
		return nil, nil
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
	case "green":
		*c = Green
	case "blue":
		*c = Blue
	default:
		return errors.New("invalid color: " + s)
	}
	return nil
}

func (c Color) String() string {
	text, err := c.MarshalText()
	if err != nil {
		return err.Error()
	}
	return string(text)
}

func Example_textMarshalUnmarshaler() {
	c := Blue

	fs := flagset.New("app")
	fs.Flag(&c, "color|c", "Color to use.")

	args := []string{"--color=green"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fs.Usage())
	fmt.Printf("Color Value: %[1]d, %[1]s\n", c)
	// Output:
	// Flags for app:
	//
	//     -c, --color  =VALUE    default: blue
	//         Color to use.
	//
	// Color Value: 1, green
}
