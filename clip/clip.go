package clip

import (
	"fmt"
    "os"
    "strings"
)

type Option interface {
	string() string
	parse(s string) error
}

type option struct {
	v           Option
	shortName   string
	longName    string
	desc        string
	hasArg      bool
	incrStep    int
	reverseFlag bool
    hide        bool
    set_        bool
}

type command struct {
	name, desc  string
	opts        []*option
	positionals []*option
	subcmds     []*command

	arguments []string
}

var RootCmd command
var cmdStack = []*command{&RootCmd}

func ArgOption(v interface{}, shortName, longName, desc string) *option {
	return RootCmd.ArgOption(v, shortName, longName, desc)
}

func ArgOptionCustom(v Option, shortName, longName, desc string) *option {
	return RootCmd.ArgOptionCustom(v, shortName, longName, desc)
}

func FlagOption(v *bool, shortName, longName, desc string) *option {
	return RootCmd.FlagOption(v, shortName, longName, desc)
}

func IncrOption(v *int, shortName, longName, desc string) *option {
	return RootCmd.IncrOption(v, shortName, longName, desc)
}

func SubCommand(name, desc string) *command {
	return RootCmd.SubCommand(name, desc)
}

func optConv(v interface{}) Option {
	var ov Option
	switch v := v.(type) {
	case *bool:
		ov = (*clipBool)(v)
	case *int:
		ov = (*clipInt)(v)
	case *int8:
		ov = (*clipInt8)(v)
	case *int16:
		ov = (*clipInt16)(v)
	case *int32:
		ov = (*clipInt32)(v)
	case *int64:
		ov = (*clipInt64)(v)
	case *uint:
		ov = (*clipUint)(v)
	case *uint8:
		ov = (*clipUint8)(v)
	case *uint16:
		ov = (*clipUint16)(v)
	case *uint32:
		ov = (*clipUint32)(v)
	case *float32:
		ov = (*clipFloat32)(v)
	case *float64:
		ov = (*clipFloat64)(v)
	case *string:
		ov = (*clipString)(v)
	default:
		panic("hello?")
	}
    return ov
}

func (c *command) Positional(v interface{}, name, desc string) *option {
    ov := optConv(v)
	o := &option{v: ov, longName: name, desc: desc}
	c.positionals = append(c.positionals, o)
	return o
}

func (c *command) PositionalCustom(v Option, name, desc string) *option {
	o := &option{v: v, longName: name, desc: desc}
	c.positionals = append(c.positionals, o)
	return o
}

func (c *command) ArgOption(v interface{}, shortName, longName, desc string) *option {
    ov := optConv(v)
	o := &option{v: ov, shortName: shortName,
		longName: longName, desc: desc, hasArg: true}
	c.opts = append(c.opts, o)
	return o
}

func (c *command) ArgOptionCustom(v Option, shortName, longName, desc string) *option {
	o := &option{v: v, shortName: shortName, desc: desc, hasArg: true}
	c.opts = append(c.opts, o)
	return o
}

func (c *command) FlagOption(v *bool, shortName, longName, desc string) *option {
	o := &option{v: (*clipBool)(v), shortName: shortName, desc: desc}
	c.opts = append(c.opts, o)
	return o
}

func (c *command) IncrOption(v *int, shortName, longName, desc string) *option {
	o := &option{v: (*clipInt)(v), shortName: shortName, desc: desc, incrStep: 1}
	c.opts = append(c.opts, o)
	return o
}

func (c *command) SubCommand(name, desc string) *command {
	sc := &command{name: name, desc: desc}
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

func PrintHelpCommand(c *command) {
    fmt.Println("Command: ", c.name)
    fmt.Println("  Options:")
    for _, o := range c.opts {
        fmt.Printf("%#v\n", o);
    }
    fmt.Println("  Positional:")
    for _, p := range c.positionals {
        fmt.Printf("%#v\n", p);
    }
    fmt.Println("  Sub-commands:")
    for _, sc := range c.subcmds {
        fmt.Printf("%s\n", sc.name);
        PrintHelpCommand(sc)
    }
}

func parseLongOpt(c *command, name string, str string) (consumed int, er error) {
    kv := strings.Split(name, "=")
    set := false
    for _, o := range c.opts {
        if kv[0] == o.longName {
            if o.hasArg {
                if len(kv) == 2 {
                    fmt.Printf("Set long option %s=%s\n", kv[0], kv[1])
                    consumed = 1
                    set, o.set_ = true, true
                } else if len(str) > 0 {
                    fmt.Printf("Set long option %s=%s\n", kv[0], str)
                    consumed = 2
                    set, o.set_ = true, true
                } else {
                    er = fmt.Errorf("optino '%s' need an argument", kv[0])
                    return
                }
            } else {
                if len(kv) > 1 {
                    er = fmt.Errorf("optino '%s' does not take argument", kv[0])
                    return
                }
                fmt.Printf("Set long option %s\n", kv[0])
                consumed = 1
                set, o.set_ = true, true
            }
        }
        if (set) {
            break
        }
    }
    if !set {
        er = fmt.Errorf("option '%s' not recognized", kv[0])
    }
    return
}

func parseShortOpt(c *command, name string, str string) (consumed int, er error) {
    for len(name) > 0 {
        var o *option
        for _, o_ := range c.opts {
            if name[:1] == o.shortName {
                o = o_
                break
            }
        }
        if o == nil {
            er = fmt.Errorf("option '%s' not recognized", name[:1])
            break
        }

        if o.hasArg {
            if len(name) > 1 {
                fmt.Printf("Set short option %s=%s\n", name[:1], name[1:])
                consumed = 1
                break
            } else if len(str) > 0 {
                fmt.Printf("Set short option %s=%s\n", name[:1], str)
                consumed = 2
                break
            } else {
                er = fmt.Errorf("option '%s' need an argument", name[:1])
                break
            }
        } else {
            fmt.Printf("Set short option %s\n", name[:1])
            name = name[1:]
            consumed = 1
        }
    }
    if er != nil {
        consumed = 0
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
            return
        } else if arg0[1] == '-' {
            if len(arg0) == 2 {
                return // '--' start arguments
            } else {
                consumed, er = parseLongOpt(c, arg0[2:], arg1)
                return
            }
        } else {
            consumed, er = parseShortOpt(c, arg0[1:], arg1)
            return
        }
    } else {
        // parsePositional or subcommand
    }
    return
}

func parseCommand(c *command, args []string) (*command, error) {
    var err error
    for len(args) > 0 {
        n := 1
        if len(args) > 1 {
            n = 2
        }
        consumed, sc, err := doParse(c, args[:n])
        if err != nil {
            c = nil
            break
        }

        if consumed > 0 {
            args = args[consumed:]
            if sc != nil {
                cmdStack = append(cmdStack, sc)
                c = sc
            }
        } else {
            c.arguments = args
            break
        }
    }
    return c, err
}

func Parse() (*command, error) {
    return parseCommand(&RootCmd, os.Args[1:])
}
