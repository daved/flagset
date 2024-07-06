package flagset_test

import (
	"fmt"

	"github.com/daved/flagset"
)

func Example() {
	var (
		help bool
		info = "default-value"
		num  int
	)

	fs := flagset.New("app")
	fs.Opt(&help, "help|h", "Display help output.")
	fs.Opt(&info, "info|i", "Info to use.")
	fs.Opt(&num, "num|n", "Number with no usage.").Meta[flagset.MetaKeySkipUsage] = true

	args := []string{"-h", "--info=non-default", "-n", "42"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Help: %t, Info: %s, Num: %d\n", help, info, num)
	fmt.Println()
	fmt.Println(fs.Usage())

	// Output:
	// Help: true, Info: non-default, Num: 42
	//
	// Flags for app:
	//
	//     -h, --help  [=BOOL]    default: false
	//         Display help output.
	//
	//     -i, --info  =STRING    default: default-value
	//         Info to use.
}

type Data struct {
	val string
}

func (d *Data) String() string {
	return d.val
}

func (d *Data) Set(s string) error {
	d.val = s
	return nil
}

func ExampleFlagValue() {
	var (
		d = &Data{val: "default-value"} // Data implements the flag.Value interface
	)

	fs := flagset.New("app")
	fs.Opt(d, "data|d", "Data to use.")

	args := []string{"--data=example"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Data: %s\n", d)
	fmt.Println()
	fmt.Println(fs.Usage())

	// Output:
	// Data: example
	//
	// Flags for app:
	//
	//     -d, --data  =VALUE    default: default-value
	//         Data to use.
}
