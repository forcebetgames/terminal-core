package domain

import (
	"fmt"
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
	fmt.Println("💰 Sistema de pagamento iniciado")
	fmt.Println("📌 Registrando hooks para teclas de dinheiro...")

	// Registra AMBAS versões: teclas normais E numpad
	p.ListenNote(2, "2", "kp_2")      // Tecla 2 normal + Numpad 2
	p.ListenNote(5, "3", "kp_3")      // Tecla 3 normal + Numpad 3
	p.ListenNote(10, "4", "kp_4")     // Tecla 4 normal + Numpad 4
	p.ListenNote(20, "5", "kp_5")     // Tecla 5 normal + Numpad 5
	p.ListenNote(50, "6", "kp_6")     // Tecla 6 normal + Numpad 6
	p.ListenNote(100, "7", "kp_7")    // Tecla 7 normal + Numpad 7

	fmt.Println("✅ Hooks de pagamento registrados!")
	fmt.Println("💡 Pressione teclas 2-7 (normais OU numpad) para inserir dinheiro")
}

func (p *PaymentCash) ListenNote(amount int, keys ...string) {
	// Register callback for all the keys
	p.InputHandler.RegisterKeyDown(keys, func(e keyboard.Event) {
		currentKey := e.Key
		currentAmount := amount

		fmt.Println("====================================")
		fmt.Printf("🔍 DEBUG - Tecla capturada:\n")
		fmt.Printf("   Tecla: %s\n", currentKey)
		fmt.Printf("   Rawcode: %d\n", e.Rawcode)
		fmt.Printf("   Valor: R$ %d\n", currentAmount)

		// Debounce - Prevenir múltiplos disparos em curto período
		p.mu.Lock()
		lastTime, exists := p.lastEventTime[currentKey]
		now := time.Now()
		if exists && now.Sub(lastTime) < 200*time.Millisecond {
			p.mu.Unlock()
			fmt.Printf("   ⏱️  REJEITADO: Debounce ativo (última tecla há %dms)\n", now.Sub(lastTime).Milliseconds())
			fmt.Println("====================================")
			return
		}
		p.lastEventTime[currentKey] = now
		p.mu.Unlock()

		fmt.Println("====================================")
		fmt.Printf("✅ VÁLIDO - Inserindo R$ %d...\n", currentAmount)

		if p.Callback != nil {
			go p.Callback(currentAmount)
		}
	})

	fmt.Printf("   ✅ Registrado: R$ %d → teclas %v\n", amount, keys)
}
