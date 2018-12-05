package clip

import (
	"fmt"
)

type ArgType interface {
    string() string
    parse(s string) error
}

type option struct {
    v ArgType
    name string
    desc string
    hasArg bool
    incrStep int
    reverseFlag bool
}

type command struct {
    opts []*option
}

var RootCmd command


func (c *command) ArgOption(v interface{}, name, desc string) *option {
    var atv ArgType
    switch v := v.(type) {
	case *int:     atv = (*clipInt)(v)
	case *int8:    atv = (*clipInt8)(v)
	case *int16:   atv = (*clipInt16)(v)
	case *int32:   atv = (*clipInt32)(v)
	case *int64:   atv = (*clipInt64)(v)
	case *uint:    atv = (*clipUint)(v)
	case *uint8:   atv = (*clipUint8)(v)
	case *uint16:  atv = (*clipUint16)(v)
	case *uint32:  atv = (*clipUint32)(v)
	case *bool:    atv = (*clipBool)(v)
    case *float32: atv = (*clipFloat32)(v)
    case *float64: atv = (*clipFloat64)(v)
    case *string:  atv = (*clipString)(v)
    default: panic("hello?")
    }
    o := &option{v:atv, name:name, desc:desc, hasArg:true}
    c.opts = append(c.opts, o)
    return o
}

func (c *command) ArgOptionCustom(v ArgType, name, desc string) *option {
    o := &option{v:v, name:name, desc:desc, hasArg:true}
    c.opts = append(c.opts, o)
    return o
}

func (c *command) FlagOption(v *bool, name, desc string) *option {
    o := &option{v:(*clipBool)(v), name:name, desc:desc, hasArg:false}
    c.opts = append(c.opts, o)
    return o
}

func (c *command) IncrOption(v *int, name, desc string) *option {
    o := &option{v:(*clipInt)(v), name:name, desc:desc}
    o.incrStep = 1 // default value
    c.opts = append(c.opts, o)
    return o
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

func ArgOption(v interface{}, name, desc string) *option {
    return RootCmd.ArgOption(v, name, desc)
}

func ArgOptionCustom(v ArgType, name, desc string) *option {
    return RootCmd.ArgOptionCustom(v, name, desc)
}

func FlagOption(v *bool, name, desc string) *option {
    return RootCmd.FlagOption(v, name, desc)
}

func IncrOption(v *int, name, desc string) *option {
    return RootCmd.IncrOption(v, name, desc)
}

func Parse(v interface{}, s string) error {
	switch v := v.(type) {
	case *int:    return ((*clipInt)(v)).parse(s)
	case *int8:   return ((*clipInt8)(v)).parse(s)
	case *int16:  return ((*clipInt16)(v)).parse(s)
	case *int32:  return ((*clipInt32)(v)).parse(s)
	case *int64:  return ((*clipInt64)(v)).parse(s)
	case *uint:   return ((*clipUint)(v)).parse(s)
	case *uint8:  return ((*clipUint8)(v)).parse(s)
	case *uint16: return ((*clipUint16)(v)).parse(s)
	case *uint32: return ((*clipUint32)(v)).parse(s)
	case *uint64: return ((*clipUint64)(v)).parse(s)
	default:      panic(fmt.Sprintf("Unsupported type %T", v))
	}
    return nil
}
