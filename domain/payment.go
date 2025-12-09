package domain

import (
	"terminal/logger"
	"sync"
	"terminal/domain/keyboard"
	"time"
)

// PaymentCash struct to handle cash payments.
type PaymentCash struct {
	Callback            func(brlCount int)
	IntervalBetweenNote int // Timeout in milliseconds between consecutive key-ups (after which the callback is called)
	InputHandler        keyboard.InputHandler
	lastEventTime       map[string]time.Time
	mu                  sync.Mutex
}

// NewPaymentCash creates a new instance of PaymentCash with the given callback.
func NewPaymentCash(callback func(brlCount int), inputHandler keyboard.InputHandler) *PaymentCash {
	return &PaymentCash{
		IntervalBetweenNote: 1000, // Timeout in milliseconds after which the callback is triggered
		Callback:            callback,
		InputHandler:        inputHandler,
		lastEventTime:       make(map[string]time.Time),
	}
}

// Start listens for key events and handles the timeout logic.
func (p *PaymentCash) Start() {
	logger.Println("💰 Sistema de pagamento iniciado")
	logger.Println("📌 Registrando hooks para teclas de dinheiro...")

	// Start listens for key events and handles the timeout logic.
	// Registra AMBAS versões: teclas normais E numpad
	p.ListenNote(2, "2", "kp_2")      // Tecla 2 normal + Numpad 2 → R$ 2
	p.ListenNote(5, "3", "kp_3")      // Tecla 3 normal + Numpad 3 → R$ 5
	p.ListenNote(10, "4", "kp_4")     // Tecla 4 normal + Numpad 4 → R$ 10
	p.ListenNote(20, "5", "kp_5")     // Tecla 5 normal + Numpad 5 → R$ 20
	p.ListenNote(50, "6", "kp_6")     // Tecla 6 normal + Numpad 6 → R$ 50
	p.ListenNote(100, "7", "kp_7")    // Tecla 7 normal + Numpad 7 → R$ 100

	logger.Println("✅ Hooks de pagamento registrados!")
	logger.Println("💡 Pressione teclas 2-7 (normais OU numpad) para inserir dinheiro")
}

func (p *PaymentCash) ListenNote(amount int, keys ...string) {
	// Register callback for all the keys
	p.InputHandler.RegisterKeyDown(keys, func(e keyboard.Event) {
		currentKey := e.Key
		currentAmount := amount

		logger.Println("====================================")
		logger.Printf("🔍 DEBUG - Tecla capturada:\n")
		logger.Printf("   Tecla: %s\n", currentKey)
		logger.Printf("   Rawcode: %d\n", e.Rawcode)
		logger.Printf("   Valor: R$ %d\n", currentAmount)

		// Debounce - Prevenir múltiplos disparos em curto período
		// Aumentado para 500ms para evitar eventos duplicados de hardware customizado
		p.mu.Lock()
		lastTime, exists := p.lastEventTime[currentKey]
		now := time.Now()
		if exists && now.Sub(lastTime) < 500*time.Millisecond {
			p.mu.Unlock()
			logger.Printf("   ⏱️  REJEITADO: Debounce ativo (última tecla há %dms)\n", now.Sub(lastTime).Milliseconds())
			logger.Println("====================================")
			return
		}
		p.lastEventTime[currentKey] = now
		p.mu.Unlock()

		logger.Println("====================================")
		logger.Printf("✅ VÁLIDO - Inserindo R$ %d...\n", currentAmount)

		if p.Callback != nil {
			go p.Callback(currentAmount)
		}
	})

	logger.Printf("   ✅ Registrado: R$ %d → teclas %v\n", amount, keys)
}
