package keyboard

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"terminal/logger"
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

	logger.Println("🎧 Iniciando escuta de eventos gohook (X11)...")

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

// mapKeycodeToKeyName maps X11 keycodes to key names
func mapKeycodeToKeyName(keycode uint16) string {
	// X11 keycode mapping
	keycodeMap := map[uint16]string{
		// Hardware buttons (custom terminal keys)
		// CORRIGIDO para ser igual ao evdev-handler.go
		// Padrão: keycode = tecla física + 1
		2: "1", // keycode 2 → tecla 1
		3: "2", // keycode 3 → tecla 2
		4: "3", // keycode 4 → tecla 3
		5: "4", // keycode 5 → tecla 4
		6: "5", // keycode 6 → tecla 5
		7: "6", // keycode 7 → tecla 6
		8: "7", // keycode 8 → tecla 7
		9: "8", // keycode 9 → tecla 8
		// Number keys (top row)
		10: "1",
		11: "2",
		12: "3",
		13: "4",
		14: "5",
		15: "6",
		16: "7",
		17: "8",
		18: "9",
		19: "r",
		// Numpad keys
		79: "kp_7", // Numpad 7/Home
		80: "kp_8", // Numpad 8/Up
		81: "kp_9", // Numpad 9/PgUp
		83: "kp_4", // Numpad 4/Left
		84: "kp_5", // Numpad 5
		85: "kp_6", // Numpad 6/Right
		87: "kp_1", // Numpad 1/End
		88: "kp_2", // Numpad 2/Down
		89: "kp_3", // Numpad 3/PgDn
		90: "kp_0", // Numpad 0/Ins
		// Letter keys
		21:  "y",    // Y key - PIX (botão físico customizado - keycode 21)
		25:  "w",
		27:  "r", // R key - Confirm withdrawal (CONFIRMAR SAQUE)
		28:  "t",
		34:  "g",    // G key - Go to home
		36:  "j",    // J key - Withdrawal modal (SACAR)
		33:  "f",    // F key - Close program
		48:  "b",    // B key - Report (RELATORIO)
		65:  " ",    // Space
		103: "up",   // Up arrow
		108: "down", // Down arrow
		// Special keys
		20:  "-",
		32:  "o",
		125: "super", // Windows/Super key
	}

	if name, ok := keycodeMap[keycode]; ok {
		return name
	}

	// Fallback to keycode number
	return fmt.Sprintf("keycode_%d", keycode)
}

// handleKeyEvent processes a key event
func (r *RealInputHandler) handleKeyEvent(ev hook.Event) {
	rawcode := int(ev.Rawcode)
	keycode := ev.Keycode
	originalKeychar := string(ev.Keychar)

	// SEMPRE usa keycode para mapear, ignorando keychar
	// Isso resolve problemas com TeamViewer que envia keychar incorreto
	keyName := mapKeycodeToKeyName(keycode)

	// Se não encontrou no mapa, tenta usar keychar como fallback
	if keyName == fmt.Sprintf("keycode_%d", keycode) {
		if originalKeychar != "" && originalKeychar != "￿" && len([]rune(originalKeychar)) > 0 {
			keyName = originalKeychar
		}
	}

	// Modo de compatibilidade TeamViewer (opcional)
	// Se TEAMVIEWER_KEYCODE_OFFSET estiver definido, aplica correção
	if offset := os.Getenv("TEAMVIEWER_KEYCODE_OFFSET"); offset != "" {
		if offsetInt, err := strconv.Atoi(offset); err == nil {
			// Aplica offset para corrigir diferença de mapeamento
			if keycode >= 10 && keycode <= 20 { // Apenas para teclas numéricas
				correctedKeycode := uint16(int(keycode) + offsetInt)
				correctedKeyName := mapKeycodeToKeyName(correctedKeycode)
				if correctedKeyName != fmt.Sprintf("keycode_%d", correctedKeycode) {
					logger.Printf("   🔧 CORREÇÃO TeamViewer: keycode %d → %d (key: '%s' → '%s')\n",
						keycode, correctedKeycode, keyName, correctedKeyName)
					keyName = correctedKeyName
					keycode = correctedKeycode
				}
			}
		}
	}

	// Debounce: prevent multiple events within 500ms for the same key
	// Aumentado para 500ms para evitar eventos duplicados de hardware customizado
	r.mu.Lock()
	lastTime, exists := r.lastEventTime[keyName]
	now := time.Now()
	if exists && now.Sub(lastTime) < 500*time.Millisecond {
		r.mu.Unlock()
		logger.Printf("⏱️  IGNORADO - Debounce: tecla '%s' (rawcode %d, keycode %d) há %dms\n",
			keyName, rawcode, keycode, now.Sub(lastTime).Milliseconds())
		return
	}
	r.lastEventTime[keyName] = now
	r.mu.Unlock()

	// Debug: show captured key with DETAILED info
	logger.Println("====================================")
	logger.Printf("⌨️  TECLA CAPTURADA (X11/gohook):\n")
	logger.Printf("   Keychar original: '%s' (bytes: %v)\n", originalKeychar, []byte(originalKeychar))
	logger.Printf("   Keycode: %d\n", keycode)
	logger.Printf("   Rawcode: %d\n", rawcode)
	logger.Printf("   Key final: '%s'\n", keyName)
	logger.Printf("   Mapeamento usado: ")
	if originalKeychar != "" && originalKeychar != "￿" && len([]rune(originalKeychar)) > 0 {
		logger.Printf("KEYCHAR direto\n")
	} else {
		logger.Printf("KEYCODE → Key (via mapa)\n")
	}
	logger.Println("====================================")

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

		logger.Printf("🔔 Executando %d callback(s) para tecla '%s'\n", len(callbacks), keyName)
		for _, callback := range callbacks {
			go callback(event)
		}
	} else {
		// Key not registered - ignore silently (optional: log for debugging)
		// fmt.Printf("⚠️  Tecla '%s' não registrada (rawcode: %d, keycode: %d)\n", keyName, rawcode, keycode)
	}
}

// SimulateKeyPress simulates a key press event (for testing purposes)
func (r *RealInputHandler) SimulateKeyPress(key string) {
	logger.Printf("⌨️  Simulando pressão de tecla: %s (gohook)\n", key)
}
