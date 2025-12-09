package keyboard

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"terminal/logger"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
)

// EvdevInputHandler handles keyboard input using evdev (works on Wayland)
type EvdevInputHandler struct {
	devices       []*evdev.InputDevice
	callbacks     map[string][]func(Event)
	mu            sync.RWMutex
	stopChan      chan struct{}
	keyCodeMap    map[uint16]string
	lastEventTime map[string]time.Time
}

// NewEvdevInputHandler creates a new evdev-based input handler
func NewEvdevInputHandler() (*EvdevInputHandler, error) {
	handler := &EvdevInputHandler{
		devices:       make([]*evdev.InputDevice, 0),
		callbacks:     make(map[string][]func(Event)),
		stopChan:      make(chan struct{}),
		lastEventTime: make(map[string]time.Time),
		keyCodeMap:    makeKeyCodeMap(),
	}

	// Find and open all keyboard devices
	if err := handler.findKeyboards(); err != nil {
		return nil, fmt.Errorf("failed to find keyboards: %w", err)
	}

	return handler, nil
}

// makeKeyCodeMap creates a mapping from evdev keycodes to our key names
func makeKeyCodeMap() map[uint16]string {
	return map[uint16]string{
		// Regular number keys (top row)
		2: "1",
		3: "2",
		4: "3",
		5: "4",
		6: "5",
		7: "6",
		8: "7",
		// Numpad keys
		79: "kp_1",
		80: "kp_2",
		81: "kp_3",
		82: "kp_4",
		83: "kp_5",
		84: "kp_6",
		85: "kp_7",
		// Command keys
		21: "y", // Y key - PIX (botão físico customizado - keycode 21)
		34: "g", // G key - Go to home
		36: "j", // J key - Withdrawal modal (SACAR)
		19: "r", // R key - Confirm withdrawal (CONFIRMAR SAQUE)
		48: "b", // B key - Report (RELATORIO)
		33: "f", // F key - Close program
		// Game controls
		13:  "=",     // Equal/Plus key - Increase bet
		12:  "-",     // Minus key - Decrease bet
		57:  " ",     // Space - Start game
		103: "up",    // Up arrow - Change game
		24:  "o",     // O key - Change game alternative
		125: "super", // Super/Windows key
	}
}

// findKeyboards finds all keyboard devices in /dev/input
func (e *EvdevInputHandler) findKeyboards() error {
	devices, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		return err
	}

	foundKeyboard := false
	for _, devicePath := range devices {
		dev, err := evdev.Open(devicePath)
		if err != nil {
			// Skip devices we can't open (permission issues, etc.)
			continue
		}

		// Check if this device looks like a keyboard based on name
		name := strings.ToLower(dev.Name)

		// EXCLUDE non-keyboard devices explicitly
		isExcluded := strings.Contains(name, "mouse") ||
			strings.Contains(name, "touchpad") ||
			strings.Contains(name, "speaker") ||
			strings.Contains(name, "lid switch") ||
			strings.Contains(name, "power button") ||
			strings.Contains(name, "sleep button") ||
			strings.Contains(name, "video bus") ||
			strings.Contains(name, "hdmi") ||
			strings.Contains(name, "headphone") ||
			strings.Contains(name, "hda") ||
			strings.Contains(name, "webcam") ||
			strings.Contains(name, "camera")

		if isExcluded {
			dev.File.Close()
			continue
		}

		// INCLUDE only devices that are clearly keyboards
		isKeyboard := strings.Contains(name, "keyboard") ||
			strings.Contains(name, "kbd") ||
			strings.Contains(name, "at translated set 2")

		if isKeyboard {
			e.devices = append(e.devices, dev)
			foundKeyboard = true
			logger.Printf("✅ Teclado detectado: %s (%s)\n", dev.Name, devicePath)
		} else {
			dev.File.Close()
		}
	}

	if !foundKeyboard {
		return fmt.Errorf("nenhum teclado encontrado. Você pode precisar executar com sudo")
	}

	return nil
}

// RegisterKeyDown registers a callback for key-down events
func (e *EvdevInputHandler) RegisterKeyDown(keys []string, callback func(Event)) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, key := range keys {
		e.callbacks[key] = append(e.callbacks[key], callback)
	}
}

// StartListening starts listening for keyboard events
func (e *EvdevInputHandler) StartListening() {
	logger.Println("🎧 Iniciando escuta de eventos evdev (Wayland)...")

	for _, dev := range e.devices {
		go e.listenDevice(dev)
	}
}

// listenDevice listens to events from a specific device
func (e *EvdevInputHandler) listenDevice(dev *evdev.InputDevice) {
	deviceName := dev.Name
	for {
		select {
		case <-e.stopChan:
			return
		default:
			events, err := dev.Read()
			if err != nil {
				log.Printf("Erro ao ler eventos de %s: %v", deviceName, err)
				return
			}

			for _, event := range events {
				// We only care about key press events (value = 1)
				// value = 0 is key release, value = 2 is key repeat
				if event.Type == evdev.EV_KEY && event.Value == 1 {
					e.handleKeyEvent(event, deviceName)
				}
			}
		}
	}
}

// handleKeyEvent processes a key event
func (e *EvdevInputHandler) handleKeyEvent(event evdev.InputEvent, deviceName string) {
	keyCode := event.Code
	keyName, ok := e.keyCodeMap[keyCode]
	if !ok {
		// Not a key we're interested in - IGNORE silently
		return
	}

	// Debounce: prevent multiple events within 500ms
	// Aumentado para 500ms para evitar eventos duplicados de hardware customizado
	e.mu.Lock()
	lastTime, exists := e.lastEventTime[keyName]
	now := time.Now()
	if exists && now.Sub(lastTime) < 500*time.Millisecond {
		e.mu.Unlock()
		logger.Printf("⏱️  IGNORADO - Debounce: tecla '%s' (keycode %d) de '%s' há %dms\n",
			keyName, keyCode, deviceName, now.Sub(lastTime).Milliseconds())
		return
	}
	e.lastEventTime[keyName] = now
	e.mu.Unlock()

	// Debug: show captured key with DETAILED info
	logger.Println("====================================")
	logger.Printf("⌨️  TECLA CAPTURADA (Wayland/evdev):\n")
	logger.Printf("   Dispositivo: %s\n", deviceName)
	logger.Printf("   Keycode evdev: %d\n", keyCode)
	logger.Printf("   Key mapeada: '%s'\n", keyName)
	logger.Printf("   Event type: %d, Value: %d\n", event.Type, event.Value)
	logger.Println("====================================")

	// Create event
	evt := Event{
		Key:     keyName,
		Rawcode: int(keyCode),
		Time:    now,
	}

	// Trigger callbacks
	e.mu.RLock()
	callbacks, exists := e.callbacks[keyName]
	e.mu.RUnlock()

	if exists {
		logger.Printf("🔔 Executando %d callback(s) para tecla '%s'\n", len(callbacks), keyName)
		for _, callback := range callbacks {
			go callback(evt)
		}
	} else {
		logger.Printf("⚠️  Nenhum callback registrado para tecla '%s'\n", keyName)
	}
}

// SimulateKeyPress simulates a key press (for testing)
func (e *EvdevInputHandler) SimulateKeyPress(key string) {
	logger.Printf("⌨️  Simulando pressão de tecla: %s (evdev)\n", key)
	// For testing purposes only
	evt := Event{
		Key:     key,
		Rawcode: 0,
		Time:    time.Now(),
	}

	e.mu.RLock()
	callbacks, exists := e.callbacks[key]
	e.mu.RUnlock()

	if exists {
		for _, callback := range callbacks {
			callback(evt)
		}
	}
}

// Stop stops listening for events
func (e *EvdevInputHandler) Stop() {
	close(e.stopChan)
	for _, dev := range e.devices {
		dev.File.Close()
	}
}

// DetectDisplayServer detects if the system is running X11 or Wayland
func DetectDisplayServer() string {
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if sessionType != "" {
		return strings.ToLower(sessionType)
	}

	// Fallback: check if WAYLAND_DISPLAY is set
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "wayland"
	}

	// Fallback: check if DISPLAY is set (X11)
	if os.Getenv("DISPLAY") != "" {
		return "x11"
	}

	// Default to x11 if we can't detect
	return "x11"
}

// NewInputHandler creates the appropriate input handler based on the display server
func NewInputHandler() (InputHandler, error) {
	displayServer := DetectDisplayServer()
	logger.Printf("🖥️  Display Server detectado: %s\n", displayServer)

	if displayServer == "wayland" {
		logger.Println("🌊 Usando modo WAYLAND (evdev)")
		logger.Println("⚠️  IMPORTANTE: Execute com sudo se houver erros de permissão!")
		handler, err := NewEvdevInputHandler()
		if err != nil {
			return nil, fmt.Errorf("falha ao criar handler evdev: %w", err)
		}
		return handler, nil
	}

	logger.Println("🪟 Usando modo X11 (gohook)")
	return NewRealInputHandler(), nil
}
