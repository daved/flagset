package flagset_test

import (
	"fmt"

	"github.com/daved/flagset"
)

func Example_flagFunc() {
	do := func(flagVal string) error {
		fmt.Println("Flag Value:", flagVal)
		return nil
	}

	fs := flagset.New("app")
	fs.Flag(do, "do|d", "Run callback.")

	args := []string{"--do=something"}

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()
	fmt.Println(fs.Usage())

	// Output:
	// Flag Value: something
	//
	// Flags for app:
	//
	//     -d, --do  =VALUE
	//         Run callback.
}
