package domain

import (
	"fmt"
	"terminal/domain/keyboard"

	hook "github.com/robotn/gohook"
)

// PaymentCash struct to handle cash payments.
type PaymentCash struct {
	Callback            func(brlCount int)
	IntervalBetweenNote int // Timeout in milliseconds between consecutive key-ups (after which the callback is called)
	InputHandler        keyboard.InputHandler
}

// NewPaymentCash creates a new instance of PaymentCash with the given callback.
func NewPaymentCash(callback func(brlCount int), inputHandler keyboard.InputHandler) *PaymentCash {
	return &PaymentCash{
		IntervalBetweenNote: 1000, // Timeout in milliseconds after which the callback is triggered
		Callback:            callback,
		InputHandler:        inputHandler,
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
	for _, keybind := range keys {
		// Captura o keybind no escopo local
		currentKey := keybind

		hook.Register(hook.KeyDown, []string{currentKey}, func(e hook.Event) {
			// Debug detalhado
			rawcode := e.Rawcode
			keychar := hook.RawcodetoKeychar(rawcode)

			fmt.Println("====================================")
			fmt.Printf("🔍 DEBUG - Tecla capturada:\n")
			fmt.Printf("   Keybind esperado: %s\n", currentKey)
			fmt.Printf("   Rawcode: %d\n", rawcode)
			fmt.Printf("   Keychar: %s\n", keychar)
			fmt.Printf("   Valor: R$ %d\n", amount)

			// Verifica se é a tecla correta (mais permissivo)
			if keychar != "" && keychar != currentKey {
				fmt.Printf("   ⚠️  Keychar diferente do esperado, mas processando mesmo assim...\n")
			}

			fmt.Println("====================================")
			fmt.Printf("💵 Inserindo R$ %d...\n", amount)

			if p.Callback != nil {
				go p.Callback(amount)
			}
		})
	}

	fmt.Printf("   ✅ Registrado: R$ %d → teclas %v\n", amount, keys)
}
