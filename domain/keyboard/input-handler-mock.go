package keyboard

import (
	"fmt"
	"time"
)

// MockInputHandler is a mock version of InputHandler for testing
type MockInputHandler struct {
	// Store key event callbacks in a map
	keyCallbacks map[string]func(Event)
}

// NewMockInputHandler creates a new instance of MockInputHandler
func NewMockInputHandler() *MockInputHandler {
	// Initialize the map for keyCallbacks
	return &MockInputHandler{
		keyCallbacks: make(map[string]func(Event)), // Ensure the map is initialized here
	}
}

// RegisterKeyDown registers a callback to be triggered on a key-down event for specific keys
func (r *MockInputHandler) RegisterKeyDown(keys []string, callback func(Event)) {
	// Register the callback for each key in the list
	for _, key := range keys {
		r.keyCallbacks[key] = callback
		fmt.Printf("Registered callback for key: %s\n", key) // Debug print to verify registration
	}
}

// SimulateKeyPress simulates a key press event (for testing purposes)
func (r *MockInputHandler) SimulateKeyPress(key string) {
	if callback, exists := r.keyCallbacks[key]; exists {
		// Simulate the event by triggering the callback
		event := Event{
			Key:     key,
			Rawcode: 0, // You can simulate any Rawcode here, or make it dynamic
			Time:    time.Now(),
		}
		callback(event)
	} else {
		// Print a message if no callback is registered for the key
		fmt.Printf("No registered callback for key: %s\n", key)
	}
}
