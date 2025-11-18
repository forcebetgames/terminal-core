package keyboard

import (
	"fmt"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
)

// RealInputHandler is responsible for handling real key press events using gohook (X11)
type RealInputHandler struct {
	started       bool
	callbacks     map[string][]func(Event)
	mu            sync.RWMutex
	lastEventTime map[string]time.Time
}

// NewRealInputHandler creates a new instance of RealInputHandler
func NewRealInputHandler() *RealInputHandler {
	return &RealInputHandler{
		started:       false,
		callbacks:     make(map[string][]func(Event)),
		lastEventTime: make(map[string]time.Time),
	}
}

// RegisterKeyDown registers a callback to be triggered on a key-down event for specific keys.
func (r *RealInputHandler) RegisterKeyDown(keys []string, callback func(Event)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, key := range keys {
		r.callbacks[key] = append(r.callbacks[key], callback)
	}
}

// StartListening starts listening for key events and triggers the registered callbacks
func (r *RealInputHandler) StartListening() {
	if r.started {
		return
	}
	r.started = true

	fmt.Println("🎧 Iniciando escuta de eventos gohook (X11)...")

	// Start gohook with a single global hook
	go func() {
		evChan := hook.Start()
		defer hook.End()

		for ev := range evChan {
			// Only process key down events
			if ev.Kind == hook.KeyDown {
				r.handleKeyEvent(ev)
			}
		}
	}()
}

// handleKeyEvent processes a key event
func (r *RealInputHandler) handleKeyEvent(ev hook.Event) {
	keyName := string(ev.Keychar)
	rawcode := int(ev.Rawcode)

	// If keychar is empty, try to use Keycode as string
	if keyName == "" {
		keyName = fmt.Sprintf("%d", ev.Keycode)
	}

	// Debounce: prevent multiple events within 200ms for the same key
	r.mu.Lock()
	lastTime, exists := r.lastEventTime[keyName]
	now := time.Now()
	if exists && now.Sub(lastTime) < 200*time.Millisecond {
		r.mu.Unlock()
		fmt.Printf("⏱️  IGNORADO - Debounce: tecla '%s' (rawcode %d) há %dms\n",
			keyName, rawcode, now.Sub(lastTime).Milliseconds())
		return
	}
	r.lastEventTime[keyName] = now
	r.mu.Unlock()

	// Debug: show captured key
	fmt.Printf("⌨️  TECLA CAPTURADA: '%s' (rawcode: %d, keycode: %d)\n", keyName, rawcode, ev.Keycode)

	// Check if we have callbacks registered for this key
	r.mu.RLock()
	callbacks, exists := r.callbacks[keyName]
	r.mu.RUnlock()

	if exists {
		event := Event{
			Key:     keyName,
			Rawcode: rawcode,
			Time:    now,
		}

		fmt.Printf("🔔 Executando %d callback(s) para tecla '%s'\n", len(callbacks), keyName)
		for _, callback := range callbacks {
			go callback(event)
		}
	} else {
		// Key not registered - ignore silently
	}
}

// SimulateKeyPress simulates a key press event (for testing purposes)
func (r *RealInputHandler) SimulateKeyPress(key string) {
	fmt.Printf("⌨️  Simulando pressão de tecla: %s (gohook)\n", key)
}
