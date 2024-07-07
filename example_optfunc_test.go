package flagset_test

import (
	"fmt"

	"github.com/daved/flagset"
)

func ExampleOptFunc() {
	do := func(flagVal string) error {
		fmt.Println("Flag Value:", flagVal)
		return nil
	}

	fs := flagset.New("app")
	fs.Opt(do, "do|d", "Run 'do'.")

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
	//         Run 'do'.
}
