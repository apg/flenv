package flenv

import (
	"strconv"
	"time"

	"net/url"
)

type value interface {
	String() string
	Set(string) error
}

// bool
type boolValue bool

func newBoolValue(p *bool) *boolValue {
	return (*boolValue)(p)
}

func (b *boolValue) Set(val string) error {
	v, err := strconv.ParseBool(val)
	*b = boolValue(v)
	return err
}

func (b *boolValue) String() string {
	if *b {
		return "true"
	}
	return "false"
}

// int
type intValue int

func newIntValue(p *int) *intValue {
	return (*intValue)(p)
}

func (i *intValue) Set(val string) error {
	v, err := strconv.Atoi(val)
	*i = intValue(v)
	return err
}

func (i *intValue) String() string {
	return strconv.Itoa(int(*i))
}

// int8
type int8Value int8

func newInt8Value(p *int8) *int8Value {
	return (*int8Value)(p)
}

func (i *int8Value) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, 8)
	*i = int8Value(v)
	return err
}

func (i *int8Value) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

// int16
type int16Value int16

func newInt16Value(p *int16) *int16Value {
	return (*int16Value)(p)
}

func (i *int16Value) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, 16)
	*i = int16Value(v)
	return err
}

func (i *int16Value) String() string {
	return strconv.FormatInt(int64(*i), 16)
}

// int32
type int32Value int32

func newInt32Value(p *int32) *int32Value {
	return (*int32Value)(p)
}

func (i *int32Value) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, 32)
	*i = int32Value(v)
	return err
}

func (i *int32Value) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

// int64
type int64Value int64

func newInt64Value(p *int64) *int64Value {
	return (*int64Value)(p)
}

func (i *int64Value) Set(val string) error {
	v, err := strconv.ParseInt(val, 10, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

// int
type uintValue uint

func newUintValue(p *uint) *uintValue {
	return (*uintValue)(p)
}

func (i *uintValue) Set(val string) error {
	v, err := strconv.Atoi(val)
	*i = uintValue(v)
	return err
}

func (i *uintValue) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

// uint8
type uint8Value uint8

func newUint8Value(p *uint8) *uint8Value {
	return (*uint8Value)(p)
}

func (i *uint8Value) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, 8)
	*i = uint8Value(v)
	return err
}

func (i *uint8Value) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

// uint16
type uint16Value uint16

func newUint16Value(p *uint16) *uint16Value {
	return (*uint16Value)(p)
}

func (i *uint16Value) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, 16)
	*i = uint16Value(v)
	return err
}

func (i *uint16Value) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

// uint32
type uint32Value uint32

func newUint32Value(p *uint32) *uint32Value {
	return (*uint32Value)(p)
}

func (i *uint32Value) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, 32)
	*i = uint32Value(v)
	return err
}

func (i *uint32Value) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

// uint64
type uint64Value uint64

func newUint64Value(p *uint64) *uint64Value {
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(val string) error {
	v, err := strconv.ParseUint(val, 10, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

// string
type stringValue string

func newStringValue(p *string) *stringValue {
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) String() string {
	return string(*s)
}

// time.Duration
type durationValue time.Duration

func newDurationValue(p *time.Duration) *durationValue {
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) String() string {
	return d.String()
}

// url.URL
type urlValue url.URL

func newURLValue(p *url.URL) *urlValue {
	return (*urlValue)(p)
}

func (u *urlValue) Set(s string) error {
	v, err := url.Parse(s)
	*u = urlValue(*v)
	return err
}

func (u *urlValue) String() string {
	return u.String()
}
