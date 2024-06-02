package flagset

import (
	"reflect"
	"testing"
)

func TestExplodeShortArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			"cmd/exp",
			[]string{"top", "-one"},
			[]string{"top", "-o", "-n", "-e"},
		},
		{
			"cmd/dub/exp/dub",
			[]string{"top", "--one", "-two", "--three"},
			[]string{"top", "--one", "-t", "-w", "-o", "--three"},
		},
		{
			"cmd/dub/exp/cmd/exp",
			[]string{"top", "--one", "-two", "sub", "-one"},
			[]string{"top", "--one", "-t", "-w", "-o", "sub", "-o", "-n", "-e"},
		},
	}

	for _, tt := range tests {
		got := explodeShortArgs(tt.args)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("\ncase: %s \n got: %s \nwant: %s ", tt.name, got, tt.want)
		}
	}
}
