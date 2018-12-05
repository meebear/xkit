package clip

import (
	"fmt"
	"strconv"
)

type (
	clipInt    int
	clipInt8   int8
	clipInt16  int16
	clipInt32  int32
	clipInt64  int64
	clipUint   uint
	clipUint8  uint8
	clipUint16 uint16
	clipUint32 uint32
	clipUint64 uint64

    clipBool   bool

    clipFloat32 float32
    clipFloat64 float64

    clipString string

    //clipDate
    //clipTime
    //clipDura
    //clipIP
)

func (i *clipInt) string() string    { return fmt.Sprintf("%d", *i) }
func (i *clipInt8) string() string   { return fmt.Sprintf("%d", *i) }
func (i *clipInt16) string() string  { return fmt.Sprintf("%d", *i) }
func (i *clipInt32) string() string  { return fmt.Sprintf("%d", *i) }
func (i *clipInt64) string() string  { return fmt.Sprintf("%d", *i) }
func (i *clipUint) string() string   { return fmt.Sprintf("%d", *i) }
func (i *clipUint8) string() string  { return fmt.Sprintf("%d", *i) }
func (i *clipUint16) string() string { return fmt.Sprintf("%d", *i) }
func (i *clipUint32) string() string { return fmt.Sprintf("%d", *i) }
func (i *clipUint64) string() string { return fmt.Sprintf("%d", *i) }

func (i *clipInt) parse(s string) (err error) {
	if s, err := strconv.ParseInt(s, 0, 0); err == nil {
		*i = clipInt(s)
	}
	return
}

func (i *clipInt8) parse(s string) (err error) {
	v, err := strconv.ParseInt(s, 0, 8)
	if err == nil {
		*i = clipInt8(v)
	}
	return
}

func (i *clipInt16) parse(s string) (err error) {
	v, err := strconv.ParseInt(s, 0, 16)
	if err == nil {
		*i = clipInt16(v)
	}
	return
}

func (i *clipInt32) parse(s string) (err error) {
	v, err := strconv.ParseInt(s, 0, 32)
	if err == nil {
		*i = clipInt32(v)
	}
	return
}

func (i *clipInt64) parse(s string) (err error) {
	v, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		*i = clipInt64(v)
	}
	return
}

func (i *clipUint) parse(s string) (err error) {
	v, err := strconv.ParseUint(s, 0, 0)
	if err == nil {
		*i = clipUint(v)
	}
	return
}

func (i *clipUint8) parse(s string) (err error) {
	v, err := strconv.ParseUint(s, 0, 8)
	if err == nil {
		*i = clipUint8(v)
	}
	return
}

func (i *clipUint16) parse(s string) (err error) {
	v, err := strconv.ParseUint(s, 0, 16)
	if err == nil {
		*i = clipUint16(v)
	}
	return
}

func (i *clipUint32) parse(s string) (err error) {
	v, err := strconv.ParseUint(s, 0, 32)
	if err == nil {
		*i = clipUint32(v)
	}
	return
}

func (i *clipUint64) parse(s string) (err error) {
	v, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		*i = clipUint64(v)
	}
	return
}

func (b *clipBool) string() string { return fmt.Sprintf("%v", *b) }

func (b *clipBool) parse(s string) (err error) {
	v, err := strconv.ParseBool(s)
	if err == nil {
		*b = clipBool(v)
	}
	return
}

func (f *clipFloat32) string() string { return fmt.Sprintf("%g", *f) }

func (f *clipFloat32) parse(s string) (err error) {
	v, err := strconv.ParseFloat(s, 32)
	if err == nil {
		*f = clipFloat32(v)
	}
	return
}

func (f *clipFloat64) string() string { return fmt.Sprintf("%g", *f) }

func (f *clipFloat64) parse(s string) (err error) {
	v, err := strconv.ParseFloat(s, 64)
	if err == nil {
		*f = clipFloat64(v)
	}
	return
}

func (s *clipString) string() string { return string(*s) }

func (s *clipString) parse(s_ string) (err error) {
    *s = clipString(s_)
	return
}
