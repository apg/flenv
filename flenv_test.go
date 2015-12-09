package flenv

import (
	"fmt"
	"log"
	"testing"
)

type config struct {
	Host string `env:"HOST,default=localhost" flag:"-h,--host" help:"Host to bind to"`
	Port int    `env:"PORT,default=80" flag:"-p,--port" help:"Port to listen on"`
}

func TestFlenvDecode(t *testing.T) {
	var cfg config
	err := Decode(&cfg)
	fmt.Printf("Error = %q\n", err)

	log.Printf("%+v\n", cfg)
}

func TestFlenvDecodeNonStruct(t *testing.T) {
	var i int
	if err := Decode(&i); err != ErrNotStruct {
		t.Fatalf("Expected error when decoding non-struct value.")
	}
}
