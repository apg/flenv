package flenv

import (
	"fmt"
	"strings"
	"testing"
)

func TestOptparseHappyWithDefaults(t *testing.T) {
	var secure bool
	var port int
	var host string

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure", false)
	opts.Add(newIntValue(&port), 'p', "port", "80", "PORT", "Port to bind on", false)
	opts.Add(newStringValue(&host), 'h', "host", "", "HOST", "Host to bind to", false)

	err := opts.Parse([]string{"-s"})
	if err != nil {
		t.Fatalf("Error not expected, got %q", err)
	}

	if !secure {
		t.Fatalf("secure should be set to true")
	}

	if port != 80 {
		t.Fatalf("port did not default to 80")
	}

	if host != "" {
		t.Fatalf("host was set")
	}
}

func TestOptparseSimpleShort(t *testing.T) {
	var secure bool

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure", false)

	err := opts.Parse([]string{"-s"})
	if err != nil {
		t.Fatalf("Error not expected, got %q", err)
	}

	if !secure {
		t.Fatalf("secure should be set to true")
	}
}

func TestOptparseSimpleLong(t *testing.T) {
	var secure bool

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure", false)

	err := opts.Parse([]string{"--secure"})
	if err != nil {
		t.Fatalf("Error not expected, got %q", err)
	}

	if !secure {
		t.Fatalf("secure should be set to true")
	}
}

func TestOptparseSimpleShortUnknown(t *testing.T) {
	var secure bool

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 0, "secure", "false", "SECURE", "Is secure", false)

	err := opts.Parse([]string{"-s"})
	if err == nil {
		t.Fatalf("Error expected, got none")
	}

	if secure {
		t.Fatalf("secure should be set to false")
	}
}

func TestOptparseSimpleShortWithArg(t *testing.T) {
	var secure bool

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure", false)

	err := opts.Parse([]string{"-s", "true"})
	if err != nil {
		t.Fatalf("Error not expected, got %q", err)
	}

	if !secure {
		t.Fatalf("secure should be set to true")
	}

	if len(opts.positions) != 0 {
		t.Fatalf("Result had positional arguments, but should not have: %d", len(opts.positions))
	}
}

func TestOptparseSimpleLongWithArg(t *testing.T) {
	var secure bool

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure", false)

	err := opts.Parse([]string{"--secure", "true"})
	if err != nil {
		t.Fatalf("Error not expected, got %q", err)
	}

	if !secure {
		t.Fatalf("secure should be set to true")
	}

	if len(opts.positions) != 0 {
		t.Fatalf("Result had positional arguments, but should not have: %d", len(opts.positions))
	}
}

func TestOptparseSimpleLongWithEqualArg(t *testing.T) {
	var secure bool

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure", false)

	err := opts.Parse([]string{"--secure=true"})
	if err != nil {
		t.Fatalf("Error not expected, got %q", err)
	}

	if !secure {
		t.Fatalf("secure should be set to true")
	}

	if len(opts.positions) != 0 {
		t.Fatalf("Result had positional arguments, but should not have: %d", len(opts.positions))
	}
}

func TestOptparsePositional(t *testing.T) {
	var secure bool

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure", false)

	err := opts.Parse([]string{"foo", "bar"})
	if err != nil {
		t.Fatalf("Error not expected, got %q", err)
	}

	if len(opts.positions) != 2 {
		t.Fatalf("Expected 2 positional arguments, had %d", len(opts.positions))
	}

	if opts.positions[0] != "foo" || opts.positions[1] != "bar" {
		t.Fatalf("Expected positional arguments to be 'foo' and 'bar', got %+v", opts.positions)
	}
}

func TestOptparsePositionalViaCancel(t *testing.T) {
	var secure bool

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure", false)

	err := opts.Parse([]string{"--", "-s"})
	if err != nil {
		t.Fatalf("Error not expected, got %q", err)
	}

	if len(opts.positions) != 1 {
		t.Fatalf("Expected 1 positional argument, had %d", len(opts.positions))
	}

	if opts.positions[0] != "-s" {
		t.Fatalf("Expected positional arguments to be '-s', got %+v", opts.positions)
	}
}

func TestOptparseShowHelp(t *testing.T) {
	var secure bool
	var port int
	var host string

	opts := OptionSet{}
	opts.Add(newBoolValue(&secure), 's', "secure", "false", "SECURE", "Is secure makes something very secure by default. This is long text that should be wordwrapped appropriately in help", false)

	opts.Add(newIntValue(&port), 'p', "port", "80", "PORT", "Port to bind on", false)
	opts.Add(newStringValue(&host), 'h', "host", "", "HOST", "Host to bind to", false)

	opts.Help()
}

func TestOptparseFillText(t *testing.T) {
	txt := `the quick brown fox`

	fmt.Println(strings.Join(fillText(txt, 10), "\n"))

	fmt.Println(strings.Join(fillText(txt, 70), "\n"))

}
