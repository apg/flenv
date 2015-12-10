package flenv

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Option struct {
	Value       value
	ShortOption byte
	LongOption  string
	Default     string // To be parsed by Value.Set() if necessary
	EnvVar      string // EnvVar to populate from if necessary
	Help        string
	Required    bool
}

func (o *Option) key() string {
	if o.ShortOption == 0 {
		return string([]byte{o.ShortOption})
	} else if o.LongOption == "" {
		return o.LongOption
	}
	return string([]byte{o.ShortOption, ','}) + o.LongOption
}

func (o *Option) Set(value string) error {
	if value == "" && o.EnvVar != "" {
		value = os.Getenv(o.EnvVar)
	}

	if value == "" {
		value = o.Default
	}

	if o.Required && value == "" {
		return fmt.Errorf("%s is required", o.key())
	}

	return o.Value.Set(value)
}

type OptionSet struct {
	shortParams map[string]*Option
	longParams  map[string]*Option
	params      map[string]*Option
	provided    map[string]*Option // option.key()
	arguments   []string
	positions   []string
}

func (s *OptionSet) Add(v value, short byte, long, defval, envvar, help string, required bool) {
	opt := &Option{
		Value:       v,
		ShortOption: short,
		LongOption:  long,
		Default:     defval,
		EnvVar:      envvar,
		Help:        help,
		Required:    required,
	}

	if short != 0 {
		if s.shortParams == nil {
			s.shortParams = make(map[string]*Option)
		}
		s.shortParams[string([]byte{short})] = opt
	}

	if long != "" {
		if s.longParams == nil {
			s.longParams = make(map[string]*Option)
		}
		s.longParams[long] = opt
	}

	if s.params == nil {
		s.params = make(map[string]*Option)
	}
	s.params[opt.key()] = opt
}

func (s *OptionSet) Set(name, value string) error {
	var opt *Option
	var ok bool
	if len(name) == 1 {
		opt, ok = s.shortParams[name]
	} else {
		opt, ok = s.longParams[name]
	}

	if !ok {
		return fmt.Errorf("option not found: %s", name)
	}

	err := opt.Set(value)
	if err != nil {
		return err
	}

	if s.provided == nil {
		s.provided = make(map[string]*Option)
	}
	s.provided[opt.key()] = opt

	return nil
}

func (s *OptionSet) Help() {
	maxParamSize := 0
	for _, opt := range s.params {
		size := 0
		if opt.LongOption != "" {
			size += len(opt.LongOption) + 2
		}
		if opt.ShortOption != 0 {
			size += 2
		}
		if opt.ShortOption != 0 && opt.LongOption != "" {
			// add 2, for ", "
			size += 2
		}

		if size > maxParamSize {
			maxParamSize = size
		}
	}

	width := 78 // account for 2 leading

	fmt.Println("Usage: ")

	for _, opt := range sortOptions(s.params) {
		// Assemble flags
		flags := ""
		if opt.ShortOption != 0 {
			flags += string([]byte{'-', opt.ShortOption})
			if opt.LongOption != "" {
				flags += ", "
			}
		}
		if opt.LongOption != "" {
			flags += "--" + opt.LongOption
		}

		// Assemble Help

		// Print it. This is wrong, as it cuts off help. Need to move to multiple lines for help.
		fmtStr := fmt.Sprintf("  %%-%ds  %%%ds\n", maxParamSize, width-maxParamSize)
		fmt.Printf(fmtStr, flags, opt.Help)
	}
}

func (s *OptionSet) Parse(args []string) error {
	var terminated bool
	s.arguments = args

	// while len(arguments), parseArg()
	for len(s.arguments) > 0 && !terminated {
		found, terminated, err := s.parseArg()
		if !found && err == nil {
			// this one must be positional
			if terminated {
				s.positions = append(s.positions, s.arguments[1:]...)
				s.arguments = s.arguments[0:0]
				break
			}

			s.positions = append(s.positions, s.arguments[0])
			if len(s.arguments) > 0 {
				s.arguments = s.arguments[1:]
			}
		} else if err != nil {
			// there's an error! So, fail somehow.
			return err
		} else if found {
			if len(s.arguments) > 0 {
				s.arguments = s.arguments[1:]
			}
		} // else parseArg advanced for us, so just continue
	}

	// trigger defaults for other arguments
	for key, opt := range s.params {
		if _, ok := s.provided[key]; !ok {
			opt.Set("") // pretend we gave nothing.
		}
	}
	return nil
}

func (s *OptionSet) parseArg() (found bool, terminated bool, err error) {
	arg := s.arguments[0]

	if len(arg) == 0 || len(arg) == 1 || (len(arg) >= 1 && arg[0] != '-') {
		return false, false, nil // blank, or '-', or positional argument
	}

	var key string
	var opt *Option
	var ok, wasEqual bool
	switch {
	case len(arg) == 2 && arg[0] == '-': // regular short arg
		if arg[1] == '-' { // end of flags, terminate
			return false, true, nil
		}

		key = arg[1:]
		opt, ok = s.shortParams[key]
		if !ok {
			return true, false, fmt.Errorf("unknown argument (%s) passed", arg)
		}

		// set us up for the value
		s.arguments = s.arguments[1:]

	case len(arg) > 2 && arg[0] == '-' && arg[1] == '-': // long arg

		idx := strings.Index(arg[2:], "=")
		if idx == 0 {
			return true, false, fmt.Errorf("invalid argument (%s) passed", arg)
		}

		key = arg[2:]
		if idx > 0 {
			key = arg[2 : idx+2]
			// We found an equal sign, so the value is *now* the remaining string.
			s.arguments[0] = arg[idx+3:]
			wasEqual = true
		}

		opt, ok = s.longParams[key]
		if !ok {
			return true, false, fmt.Errorf("unknown argument (%s) passed", key)
		}

		if !wasEqual {
			s.arguments = s.arguments[1:]
		}
	default:
		return false, false, nil
	}

	// s.arguments[0] should now be the value. Attempt to use it.
	// However, if option is boolean, argument doesn't have to be used.

	if _, ok := opt.Value.(*boolValue); ok {
		// TODO: This can likely be simplified like crazy
		if len(s.arguments) > 0 {
			err = s.Set(key, s.arguments[0])
			if err != nil && !wasEqual {
				err = s.Set(key, "true") // argument is optional
				return true, false, err
			}
			return true, false, err
		}
		err = s.Set(key, "true")
		return true, false, err
	}

	err = s.Set(key, s.arguments[0])
	return true, false, err
}

func sortOptions(opts map[string]*Option) []*Option {
	list := make(sort.StringSlice, len(opts))
	i := 0
	for _, o := range opts {
		list[i] = o.key()
		i++
	}
	list.Sort()

	result := make([]*Option, len(list))
	for i, key := range list {
		result[i] = opts[key]
	}
	return result
}
