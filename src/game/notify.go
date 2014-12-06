package game

// message is the basic type for notification messages,
// incorporating the message value and the key it was sent under.
type message struct {
	k string
	v interface{}
}

// notes keeps track of all registered channels.
var notes = make(map[string][](chan message))

// kmp stands for "keep me posted", and registers a callback channel
// for messages sent to a given key.
func kmp(key string, ch chan message) {
	notes[key] = append(notes[key], ch)
}

// notify sends a message to every channel registered for a key.
func notify(key string, value interface{}) {
	m := message{k: key, v: value}
	//fmt.Printf("sending msg %+v\n", m)
	for _, n := range notes[key] {
		n <- m
	}
}

// quit sends a nil value on every channel.
func quit() {
	notify("quit", nil)
}

// Message body types.

type locationMsg struct {
	o    interface{}
	x, y int
}
