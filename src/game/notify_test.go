package game

import (
	"fmt"
	"testing"
)

type testMessage int

func (t testMessage) String() string {
	return fmt.Sprintf("%d", int(t))
}

func TestRegister(t *testing.T) {
	kmp("TestRegister", nil)
	if got, want := len(notes["TestRegister"]), 1; got != want {
		t.Fatalf("Got len(notes[TestRegister]) = %v, want $v", got, want)
	}
	if got := notes["TestRegister"][0]; got != nil {
		t.Fatalf("Got notes[TestRegister][0] = %v, want nil", got)
	}
}

func TestNotify(t *testing.T) {
	c := make(chan message, 10)
	kmp("TestNotify", c)
	go notify("TestNotify", testMessage(5))
	m, ok := <-c
	if !ok {
		t.Fatal("Failed to read channel")
	}
	if tm, ok := m.(testMessage); ok {
		if got, want := int(tm), 5; got != want {
			t.Fatalf("Got message number %d, want %d", got, want)
		}
	} else {
		t.Fatalf("Got message of type %T, want testMessage", m)
	}
}
