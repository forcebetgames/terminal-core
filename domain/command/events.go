package command

import (
	"context"
	"os"
	"terminal/domain/keyboard"
	"terminal/logger"

	"github.com/chromedp/chromedp"
)

func DepositModal(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"y", "Y"}, func(e keyboard.Event) {
		logger.Println("Depositar / PIX pressionado (Y)")

		script := `
		window.dispatchEvent(new Event('openModalDeposit'))
		`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao abrir depósito:", errChrome)
		}
	})
}

func GotoHome(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"g"}, func(e keyboard.Event) {
		logger.Println("Ir para lista de jogos pressionado")

		script := `
				window.dispatchEvent(new Event('gotoHome'))
			`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao tentar ir para a home:", errChrome)
		}
	})
}

func ModalSaque(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"j", "J"}, func(e keyboard.Event) {
		logger.Println("abrir modal de saque (J)")

		script := `
				window.dispatchEvent(new Event('openModalWithdrawall'))
			`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao tentar abrir o modal do saque:", errChrome)
		}
	})
}

func ConfirmWithdrawal(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"r", "R"}, func(e keyboard.Event) {
		logger.Println("confirmar saque (R)")

		script := `
				window.dispatchEvent(new Event('clearAmount'))
			`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao tentar confirmar saque:", errChrome)
		}
	})
}

func ShowReport(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"b", "B"}, func(e keyboard.Event) {
		logger.Println("abrir relatório de movimentações (B)")

		script := `
				window.dispatchEvent(new Event('openReport'))
			`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao tentar abrir relatório:", errChrome)
		}
	})
}

func CloseProgram(ctx context.Context, cancel context.CancelFunc, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"f"}, func(e keyboard.Event) {
		logger.Println("Fechando terminal")

		cancel()
		os.Exit(1)
	})
}

// IncreaseBet increases the bet amount
func IncreaseBet(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"="}, func(e keyboard.Event) {
		logger.Println("Aumentar aposta pressionado (+)")

		script := `
			window.dispatchEvent(new Event('increaseBet'))
		`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao aumentar aposta:", errChrome)
		}
	})
}

// DecreaseBet decreases the bet amount
func DecreaseBet(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"-"}, func(e keyboard.Event) {
		logger.Println("Diminuir aposta pressionado (-)")

		script := `
			window.dispatchEvent(new Event('decreaseBet'))
		`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao diminuir aposta:", errChrome)
		}
	})
}

// ChangeGame changes the current game
func ChangeGame(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"up", "o", "g"}, func(e keyboard.Event) {
		logger.Printf("Mudar jogo pressionado (%s)\n", e.Key)

		script := `
			window.dispatchEvent(new Event('changeGame'))
		`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao mudar jogo:", errChrome)
		}
	})
}

// StartGame starts/spins the game
func StartGame(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{" "}, func(e keyboard.Event) {
		logger.Println("Iniciar jogo pressionado (SPACE)")

		script := `
			window.dispatchEvent(new Event('startGame'))
		`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			logger.Println("Erro ao iniciar jogo:", errChrome)
		}
	})
}
