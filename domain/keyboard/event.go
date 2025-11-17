package keyboard

import (
	"time"
)

// Event represents the event data for a key event.
type Event struct {
	Key     string
	Rawcode int
	Time    time.Time
}

// InputHandler defines methods for registering key events and handling key presses.
type InputHandler interface {
	RegisterKeyDown(keys []string, callback func(e Event))
	SimulateKeyPress(key string)
}
