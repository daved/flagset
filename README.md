# flagset [![GoDoc](https://pkg.go.dev/badge/github.com/daved/flagset.svg)](https://pkg.go.dev/github.com/daved/flagset)

```go
go get github.com/daved/flagset
```

Package flagset wraps the standard library flag package and focuses on the flag.FlagSet type. This
is done to simplify usage, and to add handling for single hyphen (e.g -h) and double hyphen (e.g.
--help) flags. Specifically, single hyphen prefixed values with multiple characters are exploded out
as though each character was its own flag (e.g. -abc = -a -b -c).

## Usage

```go
type Flag
    func (f Flag) Description() string
    func (f Flag) Longs() []string
    func (f Flag) Shorts() []string
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
}
```

## More Info

### Supported Flag Value Types

- builtin: *string, *bool, *int, *int64, *uint, *uint64, *float64
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
when the related flag is called. Both FlagFunc and FlagBoolFunc are offered so that callers can
easily associate their own functions with flags.

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

### Additional Flag Value Type Examples

[Package docs](https://pkg.go.dev/github.com/daved/flagset) contain more examples.

### Default Templating

`fs.Usage()` value from the usage example above:

```txt
Flags for app:

    -h, --help  [=BOOL]    default: false
        Display help output.

    -i, --info  =STRING    default: default-value
        Info to use.
```

### Custom Templating

Custom templates and template behaviors (i.e. template function maps) can be set. Custom data can be
attached to instances of FlagSet, and Flag using their Meta fields for access from custom templates.
