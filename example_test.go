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

	fmt.Printf("Flag Vals:\nHelp: %t, Info: %s, Num: %d\n", help, info, num)
	fmt.Println()
	fmt.Println(fs.Usage())

	// Output:
	// Flag Vals:
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
