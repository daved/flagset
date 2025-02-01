# flagset [![GoDoc](https://pkg.go.dev/badge/github.com/daved/flagset.svg)](https://pkg.go.dev/github.com/daved/flagset)

```go
go get github.com/daved/flagset
```

Package flagset provides simple flag and flag value handling using idiomatic techniques for advanced
usage. Nomenclature and handling rules are based on POSIX standards. For example, all single hyphen
prefixed arguments with multiple characters are exploded out as though they are their own flags
(e.g. -abc = -a -b -c). 

## Usage

```go
type Flag
    func (f *Flag) Description() string
    func (f *Flag) Longs() []string
    func (f *Flag) Shorts() []string
type FlagSet
    func New(name string) *FlagSet
    func (fs *FlagSet) Flag(val any, names, desc string) *Flag
    func (fs *FlagSet) Flags() []*Flag
    func (fs *FlagSet) Name() string
    func (fs *FlagSet) Operand(i int) string
    func (fs *FlagSet) Operands() []string
    func (fs *FlagSet) Parse(args []string) error
    func (fs *FlagSet) Parsed() []string
    func (fs *FlagSet) SetUsageTemplating(tmplCfg *TmplConfig)
    func (fs *FlagSet) Usage() string
// see package docs for more
```

### Setup

```go
func main() {
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
}
```

## More Info

### Supported Flag Value Types

- builtin: *string, *bool, error, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16,
*uint32, *uint64, *float32, *float64
- stdlib: *time.Duration, flag.Value
- vtype: vtype.TextMarshalUnmarshaler, vtype.FlagCallback, vtype.FlagFunc, vtype.FlagBoolFunc

#### `vtype` Types

```go
type FlagBoolFunc
    func (f FlagBoolFunc) IsBool() bool
    func (f FlagBoolFunc) OnFlag(s string) error
type FlagCallback
type FlagFunc
    func (f FlagFunc) IsBool() bool
    func (f FlagFunc) OnFlag(val string) error
type TextMarshalUnmarshaler
```

The main vtype types are interface types. First, TextMarshalUnmarshaler describes types which
satisfy both the encoding.TextMarshaler and encoding.TextUnmarshaler interfaces, and is offered so
that callers can easily use standard library compatible types. Second, FlagCallback describes types
which indicate whether they are intended to be used with bool flags and provide an action to take
when the related flag is called. Both FlagFunc and FlagBoolFunc implement FlagCallback and are
offered so that callers can easily associate their own functions with flags. That is, compatible
functions will be automatically converted to either FlagFunc or FlagBoolFunc.

```go
func main() {
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
}
```
Output:
```txt
Flag Value: something
```

### Additional Flag Value Type Examples

[Package docs](https://pkg.go.dev/github.com/daved/flagset) contain more examples.

### Default Templating

`fs.Usage()` value from the usage example above:

```txt
Flags for app:

    -i, --info  =STRING    default: default-value
        Interesting info.

    -v, --verbose  [=BOOL]    default: false
        Set verbose output.
```

### Custom Templating

Custom templates and template behaviors (i.e. template function maps) can be set. Custom data can be
attached to instances of FlagSet, and Flag using their Meta fields for access from custom templates.
