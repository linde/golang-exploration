package greeterserver

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestGreeterServerInstance(t *testing.T) {
	greeterserver := New()
	if greeterserver == nil {
		t.Fatalf(`New() is nil`)
	}
}
