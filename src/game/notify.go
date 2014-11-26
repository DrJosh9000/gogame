package game

import (
	"fmt"
)

// message is the basic type for notification messages,
// incorporating the message value and the key it was sent under.
type message struct {
	k string
	v fmt.Stringer
}

func (m *message) String() string {
	return fmt.Sprintf("message{k:%s, v:%+v}", m.k, m.v)
}

// notes keeps track of all registered channels.
var notes = map[string][](chan message){}

// kmp stands for "keep me posted", and registers a callback channel
// for messages sent to a given key.
func kmp(key string, ch chan message) {
	notes[key] = append(notes[key], ch)
}

// notify sends a message to every channel registered for a key.
func notify(key string, value fmt.Stringer) {
	m := message{k: key, v: value}
	//fmt.Printf("sending msg %+v\n", m)
	for _, n := range notes[key] {
		n <- m
	}
}

// Message body types.

type basicMsg string

func (m basicMsg) String() string {
	return string(m)
}

var (
	quitMsg = basicMsg("quit")
)

type locationMsg struct {
	o    fmt.Stringer
	x, y int
}

func (l locationMsg) String() string {
	return fmt.Sprintf("locationMsg{o:%v x:%d y:%d}", l.o, l.x, l.y)
}
