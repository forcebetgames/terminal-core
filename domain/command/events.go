package command

import (
	"context"
	"fmt"
	"os"
	"terminal/domain/keyboard"

	"github.com/chromedp/chromedp"
)

func DepositModal(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"y"}, func(e keyboard.Event) {
		fmt.Println("Depositar pressionado")

		script := `
		window.dispatchEvent(new Event('openModalDeposit'))
		`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			fmt.Println("Erro ao abrir depósito:", errChrome)
		}
	})
}

func GotoHome(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"g"}, func(e keyboard.Event) {
		fmt.Println("Ir para lista de jogos pressionado")

		script := `
				window.dispatchEvent(new Event('gotoHome'))
			`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			fmt.Println("Erro ao tentar ir para a home:", errChrome)
		}
	})
}

func ModalSaque(ctx context.Context, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"j"}, func(e keyboard.Event) {
		fmt.Println("abrir modal de saque")

		script := `
				window.dispatchEvent(new Event('openModalWithdrawall'))
			`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil),
		)
		if errChrome != nil {
			fmt.Println("Erro ao tentar abrir o modal do saque:", errChrome)
		}
	})
}

func CloseProgram(ctx context.Context, cancel context.CancelFunc, inputHandler keyboard.InputHandler) {
	inputHandler.RegisterKeyDown([]string{"f"}, func(e keyboard.Event) {
		fmt.Println("Fechando terminal")

		cancel()
		os.Exit(1)
	})
}
