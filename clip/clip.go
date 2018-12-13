package clip

import (
    "fmt"
    "os"
    "time"
    "strings"
    "net"
    "bytes"
)

// need pointer receiver
type Option interface {
    String() string
    Parse(s string) error
}

type optSt int
const (
    optStDefault optSt = iota
    optStMustSet
    optStSet
);

type option struct {
    v           Option
    shortName   byte
    longName    string
    desc        string

    hasArg      bool
    incrStep    int

    reverseFlag bool
    hide        bool
    repeatable  bool

    status      optSt
}

type command struct {
    name, desc  string
    opts        []*option
    positionals []*option
    subcmds     []*command

    Arguments   []string

    parent_     *command
}

var helpOption = option{ shortName: 'h', longName: "help",
                    desc: "command usage" }

var RootCmd command
var progInfo string

func ArgOption(v interface{}, shortName byte, longName, desc string) *option {
    return RootCmd.ArgOption(v, shortName, longName, desc)
}

func ArgOptionCustom(v Option, shortName byte, longName, desc string) *option {
    return RootCmd.ArgOptionCustom(v, shortName, longName, desc)
}

func FlagOption(v *bool, shortName byte, longName, desc string) *option {
    return RootCmd.FlagOption(v, shortName, longName, desc)
}

func IncrOption(v *int, shortName byte, longName, desc string) *option {
    return RootCmd.IncrOption(v, shortName, longName, desc)
}

func Positional(v interface{}, name, desc string) *option {
    return RootCmd.Positional(v, name, desc)
}

func SubCommand(name, desc string) *command {
    return RootCmd.SubCommand(name, desc)
}

func optConv(v interface{}) Option {
    var ov Option
    switch v := v.(type) {
    case *bool:   ov = (*clipBool)(v)
    case *int:    ov = (*clipInt)(v)
    case *int8:   ov = (*clipInt8)(v)
    case *int16:  ov = (*clipInt16)(v)
    case *int32:  ov = (*clipInt32)(v)
    case *int64:  ov = (*clipInt64)(v)
    case *uint:   ov = (*clipUint)(v)
    case *uint8:  ov = (*clipUint8)(v)
    case *uint16: ov = (*clipUint16)(v)
    case *uint32: ov = (*clipUint32)(v)
    case *float32:ov = (*clipFloat32)(v)
    case *float64:ov = (*clipFloat64)(v)
    case *string: ov = (*clipString)(v)
    case *time.Duration: ov = (*clipDura)(v)
    case *net.IP: ov = (*clipIP)(v)
    default: panic(fmt.Sprintf("use _Custom() for option type %T", v))
    }
    return ov
}

func (c *command) Positional(v interface{}, name, desc string) *option {
    o := &option{v: optConv(v), longName: name, desc: desc}
    c.positionals = append(c.positionals, o)
    return o
}

func (c *command) PositionalCustom(v Option, name, desc string) *option {
    o := &option{v: v, longName: name, desc: desc}
    c.positionals = append(c.positionals, o)
    return o
}

func (c *command) appendOption(o *option) *option {
    if o.shortName == helpOption.shortName {
        helpOption.shortName = 0
    }
    if o.longName == helpOption.longName {
        helpOption.longName = ""
    }
    c.opts = append(c.opts, o)
    return o
}

func (c *command) ArgOption(v interface{}, shortName byte, longName, desc string) *option {
    o := &option{v: optConv(v), shortName: shortName, longName: longName,
                 desc: desc, hasArg: true}
    return c.appendOption(o)
}

func (c *command) ArgOptionCustom(v Option, shortName byte, longName, desc string) *option {
    o := &option{v: v, shortName: shortName, desc: desc, hasArg: true}
    return c.appendOption(o)
}

func (c *command) FlagOption(v *bool, shortName byte, longName, desc string) *option {
    o := &option{v: (*clipBool)(v), shortName: shortName, longName: longName, desc: desc}
    return c.appendOption(o)
}

func (c *command) IncrOption(v *int, shortName byte, longName, desc string) *option {
    o := &option{v: (*clipInt)(v), shortName: shortName, longName: longName,
        desc: desc, incrStep: 1, repeatable: true}
    return c.appendOption(o)
}

func SetHelpOption(shortName byte, longName string) {
    helpOption.shortName = shortName
    helpOption.longName = longName
}

func (c *command) SubCommand(name, desc string) *command {
    sc := &command{name: name, desc: desc, parent_: c}
    c.subcmds = append(c.subcmds, sc)
    return sc
}

func (o *option) SetIncrStep(step int) *option {
    if o.incrStep == 0 {
        panic("cannot set increment step on non-increment option")
    }
    if step == 0 {
        panic("increment step cannot be 0")
    }
    o.incrStep = step
    return o
}

func (o *option) ReverseFlag() *option {
    if _, ok := o.v.(*clipBool); !ok {
        panic("ReverseFlag on non-bool option")
    }
    o.reverseFlag = true
    return o
}

func (o *option) Hide() *option {
    o.hide = true
    return o
}

func (o *option) Repeatable(r bool) *option {
    o.repeatable = r
    return o
}

func (o *option) MustSet() *option {
    o.status = optStMustSet
    return o
}

func errf(format string, args ...interface{}) error {
    return fmt.Errorf(fmt.Sprintf("CommandLine: %s", format), args...)
}

func prtf(format string, args ...interface{}) {
    fmt.Printf(fmt.Sprintf("CommandLine: %s", format), args...)
}

func setNoArgOption(o *option) {
    if o.incrStep != 0 {
        if v_, ok := o.v.(*clipInt); ok {
            v := (*int)(v_);
            *v += o.incrStep
        } else {
            panic("internal: none integer option has non-zero incrStep")
        }
    } else {
        if v_, ok := o.v.(*clipBool); ok {
            v := (*bool)(v_)
            *v = !o.reverseFlag
        } else {
            panic("internal: none boolean option has zero incrStep")
        }
    }
}

func parseLongOpt(c *command, name string, str string) (consumed int, er error) {
    kv := strings.Split(name, "=")
    set := false
    for _, o := range c.opts {
        if o == &helpOption {
            continue
        }
        if kv[0] == o.longName {
            if o.status == optStSet && !o.repeatable {
                er = errf("option '%s' set more than once", kv[0])
                return
            }
            if o.hasArg {
                if len(kv) == 2 {
                    if er = o.v.Parse(kv[1]); er != nil { return }
                    prtf("Set long option %s=%s\n", kv[0], kv[1])
                    consumed = 1
                    set, o.status = true, optStSet
                } else if len(str) > 0 {
                    if er = o.v.Parse(str); er != nil { return }
                    prtf("Set long option %s=%s\n", kv[0], str)
                    consumed = 2
                    set, o.status = true, optStSet
                } else {
                    er = errf("optino '%s' need an argument", kv[0])
                    return
                }
            } else {
                if len(kv) > 1 {
                    er = errf("optino '%s' does not take argument", kv[0])
                    return
                }
                setNoArgOption(o)
                prtf("Set long option %s\n", kv[0])
                consumed = 1
                set, o.status = true, optStSet
            }
        }
        if (set) {
            break
        }
    }

    if !set {
        if kv[0] == helpOption.longName {
            HelpCommand(c)
            os.Exit(0)
        }
        consumed = 0
        if er == nil {
            er = errf("option '%s' not recognized", kv[0])
        }
    }
    return
}

func parseShortOpt(c *command, name string, str string) (consumed int, er error) {
    for len(name) > 0 {
        var o *option
        for _, o_ := range c.opts {
            if o_ == &helpOption {
                continue
            }
            if name[0] == o_.shortName {
                o = o_
                break
            }
        }
        if o == nil || o.v == nil {
            if name[0] == helpOption.shortName {
                HelpCommand(c)
                os.Exit(0)
            }
            er = errf("option '%s' not recognized", name[:1])
            break
        }
        if o.status == optStSet && !o.repeatable {
            er = errf("optino '%s' set more than once", name[:1])
            break
        }

        if o.hasArg {
            if len(name) > 1 {
                if er = o.v.Parse(name[1:]); er != nil { return }
                prtf("Set short option %s=%s\n", name[:1], name[1:])
                consumed = 1
                o.status = optStSet
                break
            } else if len(str) > 0 {
                if er = o.v.Parse(str); er != nil { return }
                prtf("Set short option %s=%s\n", name[:1], str)
                consumed = 2
                o.status =  optStSet
                break
            } else {
                er = errf("option '%s' need an argument", name[:1])
                break
            }
        } else {
            setNoArgOption(o)
            prtf("Set short option %s\n", name[:1])
            name = name[1:]
            consumed = 1
            o.status = optStSet
        }
    }
    if er != nil {
        consumed = 0
    }
    return
}

func parsePositional(c *command, str string) (consumed int, er error) {
    for _, o := range c.positionals {
        if o.status == optStSet {
            continue
        }
        if er = o.v.Parse(str); er != nil { return }
        prtf("Set positianl '%s' to '%s'\n", o.longName, str)
        o.status = optStSet
        consumed = 1
        break
    }
    return
}

func parseSubCommand(c *command, str string) (consumed int, sc *command, er error) {
    for _, s := range c.subcmds {
        if s.name == str {
            sc = s
            consumed = 1
            break
        }
    }
    return
}

func doParse(c *command, ss []string) (consumed int, sc *command, er error) {
    arg0 := ss[0]
    var arg1 string
    if len(ss) > 1 {
        arg1 = ss[1]
    }

    if arg0[0] == '-' {
        if len(arg0) == 1 {
            fmt.Println("warning: option '-' ignored")
            consumed = 1
        } else if arg0[1] == '-' {
            if len(arg0) > 2 {
                consumed, er = parseLongOpt(c, arg0[2:], arg1)
            }
        } else {
            consumed, er = parseShortOpt(c, arg0[1:], arg1)
        }
    } else {
        if consumed, er = parsePositional(c, arg0); er == nil {
            if consumed == 0 {
                consumed, sc, er = parseSubCommand(c, arg0)
            }
        }
    }
    return
}

func checkMustSetOptions(c *command) error {
    for c != nil {
        for _, o := range c.opts {
            if o.status == optStMustSet {
                return fmt.Errorf("option '%s' not given", o.longName) //fixme
            }
        }
        for _, o := range c.positionals {
            if o.status == optStMustSet {
                return fmt.Errorf("positional option '%s' not given", o.longName)
            }
        }
        c = c.parent_
    }
    return nil
}

func parseCommand(c *command, args []string) (*command, error) {
    var err error
    for len(args) > 0 {
        n := 1
        if len(args) > 1 {
            n = 2
        }
        consumed, sc, er := doParse(c, args[:n])
        if er != nil {
            err = er
            c = nil
            break
        }

        if consumed > 0 {
            args = args[consumed:]
            if sc != nil {
                c = sc
            }
        } else {
            c.Arguments = args
            break
        }
    }
    if err = checkMustSetOptions(c); err != nil {
        c = nil
    }
    return c, err
}

func Parse() (*command, error) {
    if helpOption.shortName != 0 || len(helpOption.longName) > 0 {
        RootCmd.opts = append(RootCmd.opts, &helpOption)
    }
    return parseCommand(&RootCmd, os.Args[1:])
}

func formatText(text string, width, indent, indentFrom uint) string {
    var buf bytes.Buffer
    indstr := "\n"

    if indent > 0 {
        buf.WriteByte('\n')
        for i:=0; i<int(indent); i++ {
            buf.WriteByte(' ')
        }
        indstr = buf.String()
        buf.Reset()

        if indentFrom == 0 {
            buf.Write([]byte(indstr[1:]))
        }
    }

    var w, wlen int
    var word string
    for len(text) > 0 {
        if ix := strings.IndexAny(text, " "); ix >= 0 {
            wlen = ix + 1
            word = text[:wlen]
            text = text[wlen:]
        } else {
            word = text
            text = ""
        }

        if w + wlen > int(width) + 1 {
            buf.Write([]byte(indstr))
            w = wlen
        } else {
            w += wlen
        }
        buf.Write([]byte(word))
    }

    return buf.String()
}

func prtList(lst [][2]string, kind string) (n int) {
    var w int
    for _, e := range lst {
        if w < len(e[0]) && len(e[0]) < 32 {
            w = len(e[0])
        }
    }

    if w < 20 { w = 20 }
    if w > 32 { w = 32 }
    w += 2
    for i, o := range lst {
        if i == 0 {
            fmt.Printf("%s:\n\n", kind)
        }
        if len(o[0]) > w-2 {
            fmt.Printf("%s\n", o[0])
            fmt.Printf("%s\n", formatText(o[1], uint(80-w), uint(w), 0))
        } else {
            fmt.Printf("%-[1]*s", w, o[0])
            fmt.Printf("%s\n", formatText(o[1], uint(80-w), uint(w), 1))
        }
        n++
    }
    return n
}

func prtOptions(os []*option, kind string) {
    var buf bytes.Buffer
    var lst [][2]string

    for _, o := range os {
        buf.Reset()
        buf.Write([]byte("  "))
        if o.shortName != 0 {
            buf.WriteByte('-')
            buf.WriteByte(o.shortName)
            buf.WriteByte(',')
        }
        if len(o.longName) > 0 {
            buf.Write([]byte(fmt.Sprintf("--%s", o.longName)))
        }
        ostr := buf.String()

        buf.Reset()
        buf.Write([]byte(o.desc))
        if o.v != nil {
            if o.status == optStDefault {
                dft := o.v.String()
                if len(dft) > 0 {
                    buf.Write([]byte(fmt.Sprintf(" (default: %s)", dft)))
                }
            } else if o.status == optStMustSet {
                buf.Write([]byte(" (must set)"))
            }
        }

        desc := buf.String()
        lst = append(lst, [2]string{ostr, desc})
    }

    if prtList(lst, kind) > 0 {
        fmt.Println()
    }
}

func HelpCommand(c *command) {
    if (c != &RootCmd) {
        fmt.Printf("Command: %s\n", c.name)
    }
    prtOptions(c.opts, "Options")
    prtOptions(c.positionals, "Positionals")

    var lst [][2]string
    for _, sc := range c.subcmds {
        lst = append(lst, [2]string{fmt.Sprintf("  %s", sc.name), sc.desc})
    }
    if prtList(lst, "Sub-Commands") > 0 {
        fmt.Println()
    }
}

func Help() {
    fmt.Printf("%s\n", formatText(progInfo, 80, 0, 0))
    HelpCommand(&RootCmd)
}

func ProgDescription(desc string) {
    progInfo = desc
}
