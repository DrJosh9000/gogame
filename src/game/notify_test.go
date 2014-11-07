package game

import (
	"fmt"
	"testing"
)

type testMessage int
func (t testMessage) String() string {
	return fmt.Sprintf("%d", int(t))
}

func TestNotify(t *testing.T) {
	c := make(chan message)
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