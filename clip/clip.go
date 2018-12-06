package clip

import (
	"fmt"
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
	mustHasArg  bool
	incrStep    int
	reverseFlag bool
}

type command struct {
	name, desc  string
	opts        []*option
	positionals []*option
	subcmds     []*command

	arguments []string
}

var RootCmd command

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
		longName: longName, desc: desc, mustHasArg: true}
	c.opts = append(c.opts, o)
	return o
}

func (c *command) ArgOptionCustom(v Option, shortName, longName, desc string) *option {
	o := &option{v: v, shortName: shortName, desc: desc, mustHasArg: true}
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
