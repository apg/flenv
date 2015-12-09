package flenv

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ErrNotStruct is an error
var ErrNotStruct = errors.New("Attempted to decode into a value that was not a pointer to struct")

// ErrNotStruct is an error
var ErrUnsupportedFlagType = errors.New("Unsupported flag type")

// DecodeArgs runs decode with os.Args
func DecodeArgs(result interface{}) (*flag.FlagSet, error) {
	return Decode(result, os.Args)
}

// Decode attempts to populate result, based on struct tags from args and the environment
//
// `result` should be a struct with at least one of the following tags
//
//    env:"ENVIRONMENT_VARIABLE"
//    flag:"--environment-variable"
//
// Additionally, the following tags are also allowed
//
//    help:"help text for flag parsing"
//    default:"default value"
func Decode(result interface{}, args []string) (*flag.FlagSet, error) {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr {
		return nil, ErrNotStruct
	}

	st := val.Elem()
	if st.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	typ := val.Elem().Type()

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		typeField := typ.Field(i)

		// Get env tag.
		envName := typeField.Tag.Get("env")
		defaultValue := typeField.Tag.Get("default")

		// Set the value to the default value
		if err := setValue(field, envName, defaultValue); err != nil {
			return nil, err
		}

		// If there's a flag associated, we need to set that up, possibly with a default.
		flagNames := decodeFlagTag(typeField.Tag.Get("flag"))
		flagHelp := typeField.Tag.Get("help")
		if len(flagNames) > 0 {
			addFlag(flagSet, field, flagNames, flagHelp)
		}
	}

	return flagSet, flagSet.Parse(args)
}

func decodeFlagTag(tag string) (names []string) {
	bits := strings.Split(tag, ",")
	out := make([]string, 0, len(bits))
	for _, bit := range bits {
		trimmed := strings.TrimLeft(bit, "-")
		if len(trimmed) > 0 {
			out = append(out, trimmed)
		}
	}

	return out
}

func addFlag(flagSet *flag.FlagSet, field reflect.Value, names []string, help string) error {
	for _, name := range names {
		switch field.Kind() {
		case reflect.Bool:
			flagSet.BoolVar(field.Addr().Interface().(*bool), name, field.Bool(), help)
		case reflect.Float32, reflect.Float64:
			flagSet.Float64Var(field.Addr().Interface().(*float64), name, field.Float(), help)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			flagSet.IntVar(field.Addr().Interface().(*int), name, int(field.Int()), help)
		case reflect.Int64:
			if t := field.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
				flagSet.DurationVar(field.Addr().Interface().(*time.Duration), name, time.Duration(field.Int()), help)
			} else {
				flagSet.Int64Var(field.Addr().Interface().(*int64), name, field.Int(), help)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			flagSet.Uint64Var(field.Addr().Interface().(*uint64), name, field.Uint(), help)
		case reflect.String:
			flagSet.StringVar(field.Addr().Interface().(*string), name, field.String(), help)
		default:
			// TODO: We can make a URL parser, but not right now.
			return ErrUnsupportedFlagType
		}
	}
	return nil
}

func setValue(field reflect.Value, name, def string) error {
	tmp := os.Getenv(name)
	if tmp == "" {
		tmp = def
	}

	if tmp == "" {
		return nil
	}

	switch field.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(tmp)
		if err != nil {
			return err
		}
		field.SetBool(v)

	case reflect.Float32, reflect.Float64:
		bits := field.Type().Bits()
		v, err := strconv.ParseFloat(tmp, bits)
		if err != nil {
			return err
		}
		field.SetFloat(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t := field.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
			v, err := time.ParseDuration(tmp)
			if err != nil {
				return err
			}
			field.SetInt(int64(v))

		} else {
			bits := field.Type().Bits()
			v, err := strconv.ParseInt(tmp, 0, bits)
			if err != nil {
				return err
			}
			field.SetInt(v)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bits := field.Type().Bits()
		v, err := strconv.ParseUint(tmp, 0, bits)
		if err != nil {
			return err
		}
		field.SetUint(v)

	case reflect.String:
		field.SetString(tmp)

	case reflect.Ptr:
		if t := field.Type().Elem(); t.Kind() == reflect.Struct && t.PkgPath() == "net/url" && t.Name() == "URL" {
			v, err := url.Parse(tmp)
			if err != nil {
				return err
			}

			field.Set(reflect.ValueOf(v))
		}
	}
	return nil
}
