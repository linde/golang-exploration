package greeterserver

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestGreeterServerInstance(t *testing.T) {
	gs := GreeterServerInstance()
	if gs == nil {
		t.Fatalf(`GreeterServerInstance() is nil`)
	}
}
