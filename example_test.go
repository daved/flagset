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
	fs.Flag(&help, "help|h", "Display help output.")
	fs.Flag(&info, "info|i", "Info to use.")
	fs.Flag(&num, "num|n", "Number with no usage.").HideUsage = true

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
