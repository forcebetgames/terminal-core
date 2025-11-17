package keyboard

import (
	"fmt"
	"time"

	hook "github.com/robotn/gohook"
)

// RealInputHandler is responsible for handling real key press events using gohook
type RealInputHandler struct{}

// NewRealInputHandler creates a new instance of RealInputHandler
func NewRealInputHandler() *RealInputHandler {
	return &RealInputHandler{}
}

// RegisterKeyDown registers a callback to be triggered on a key-down event for specific keys.
func (r *RealInputHandler) RegisterKeyDown(keys []string, callback func(Event)) {
	// Register the keys using gohook
	for _, key := range keys {
		hook.Register(hook.KeyUp, []string{""}, func(e hook.Event) {
			if hook.RawcodetoKeychar(e.Rawcode) != key {
				return
			}
			event := Event{
				Key:     key,
				Rawcode: int(e.Rawcode),
				Time:    time.Now(),
			}
			callback(event)
		})
	}
}

// StartListening starts listening for key events and triggers the registered callbacks
func (r *RealInputHandler) StartListening() {
	// Start listening for events
	hook.Start()
}

// SimulateKeyPress simulates a key press event (for testing purposes)
func (r *RealInputHandler) SimulateKeyPress(key string) {
	// Simulate a key press event (this part would be used for testing or simulations)
	fmt.Printf("Simulated key press: %s\n", key)
}
