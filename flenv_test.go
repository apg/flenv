package flenv

import (
	"os"
	"testing"
)

type config struct {
	Host string `env:"HOST" default:"localhost" flag:"-h,--host" help:"Host to bind to"`
	Port int    `env:"PORT" default:"80" flag:"-p,--port" help:"Port to listen on"`
}

func TestFlenvDecode(t *testing.T) {
	oldHost := os.Getenv("HOST")
	oldPort := os.Getenv("PORT")

	os.Setenv("HOST", "sigusr2.net")
	os.Setenv("PORT", "8000")

	defer func() {
		os.Setenv("HOST", oldHost)
		os.Setenv("PORT", oldPort)
	}()

	var cfg config
	_, err := Decode(&cfg, []string{})
	if err != nil {
		t.Fatalf("Unexpected error, got %q", err)
	}

	if cfg.Host != "sigusr2.net" {
		t.Fatalf("Expected host to be \"sigusr2.net\", got %q", cfg.Host)
	}
	if cfg.Port != 8000 {
		t.Fatalf("Expected port to be \"8000\", got %q", cfg.Port)
	}
}

func TestFlenvDecodeArgs(t *testing.T) {
	var cfg config
	fs, err := Decode(&cfg, []string{"-h", "sigusr2.net", "-p", "8000"})
	if err != nil {
		t.Fatalf("Unexpected error, got %q", err)
	}
	if fs == nil {
		t.Fatalf("Flag set was nil")
	}
	if cfg.Host != "sigusr2.net" {
		t.Fatalf("Expected host to be \"sigusr2.net\", got %q", cfg.Host)
	}
	if cfg.Port != 8000 {
		t.Fatalf("Expected port to be \"8000\", got %q", cfg.Port)
	}
}

func TestFlenvDecodeArgsWithEnv(t *testing.T) {
	oldHost := os.Getenv("HOST")
	oldPort := os.Getenv("PORT")

	os.Setenv("HOST", "sigusr2.net")
	os.Setenv("PORT", "8000")

	defer func() {
		os.Setenv("HOST", oldHost)
		os.Setenv("PORT", oldPort)
	}()

	var cfg config
	_, err := Decode(&cfg, []string{})
	if err != nil {
		t.Fatalf("Unexpected error, got %q", err)
	}
	if cfg.Host != "sigusr2.net" {
		t.Fatalf("Expected host to be \"sigusr2.net\", got %q", cfg.Host)
	}
	if cfg.Port != 8000 {
		t.Fatalf("Expected port to be \"8000\", got %q", cfg.Port)
	}
}

func TestFlenvDecodeNonStruct(t *testing.T) {
	var i int
	if _, err := Decode(&i, []string{}); err != ErrNotStruct {
		t.Fatalf("Expected error when decoding non-struct value.")
	}
}
