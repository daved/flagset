package flagset_test

import (
	"fmt"
	"net/url"

	"github.com/daved/flagset"
)

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
	if u, err := url.Parse(s); err != nil {
		return err
	} else {
		*v.URL = *u
	}
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

	fmt.Printf("URL: %s\n", url)
	fmt.Println()
	fmt.Println(fs.Usage())

	// Output:
	// URL: https://example.com
	//
	// Flags for app:
	//
	//     -u, --url  =VALUE
	//         URL to use.
}
