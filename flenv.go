package flenv

import (
	"errors"
	"log"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ErrNotStruct = errors.New("Attempted to decode into a value that was not a pointer to struct")

func Decode(result interface{}) error {
	//	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr {
		return ErrNotStruct
	}

	st := val.Elem()
	if st.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	typ := val.Elem().Type()

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		typeField := typ.Field(i)

		log.Printf("%+v\n", typ)

		var defaultValue string

		// Get env tag.
		envName, envDef := decodeEnvTag(typeField.Tag.Get("env"))
		if envDef != "" {
			defaultValue = envDef
		}

		log.Printf("Field=%s has default = %s, but will be populated from %s\n", typeField.Name, defaultValue, envName)
		log.Printf("Field: %+v            TypeField: %+v\n", field, typeField)

		if err := setValue(field, envName, envDef); err != nil {
			return err
		}

		/// TODO: Deal with flags.
		// OK, now we *may* have a way to set the value for this
		// field. But, even if we do, it may be overridden by a flag.
		// flagNames := decodeFlagTag(typeField.Tag.Get("flag"))
		// flagHelp := typeField.Tag.Get("help")

		//  if len(flagNames) > 0 {
		//	}
	}

	return nil
}

// Decodes the name given in tag.
func decodeEnvTag(tag string) (name string, def string) {
	idx := strings.Index(tag, ",")
	if idx > 0 {
		name = tag[0:idx]
	}

	idx = strings.Index(tag, "default=")
	if idx > 0 {
		def = tag[idx+8:]
	}

	return
}

func setValue(field reflect.Value, name, def string) error {
	tmp := os.Getenv(name)
	if tmp == "" {
		tmp = def
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
