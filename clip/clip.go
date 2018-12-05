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
}

type command struct {
    opts []option
}

var RootCmd command

func (c *command) ArgOption(v ArgType, name, desc string) *command {
    return c
}

func ArgOption(v ArgType, name, desc string) *command {
    return &RootCmd
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
