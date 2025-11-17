package command

import (
	"context"
	"fmt"
	"os"

	"github.com/chromedp/chromedp"

	hook "github.com/robotn/gohook"
)

func init() {
	// Define key codes for Linux/Windows/MacOS
	// if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
	// 	KeyJ = hotkey.Key(0x24) // Linux/Windows specific key code for 'J'
	// 	KeyG = hotkey.Key(0x22) // Linux/Windows specific key code for 'G'
	// } else {
	// 	// Default for macOS
	// 	KeyJ = hotkey.KeyJ
	// 	KeyG = hotkey.KeyG
	// }
}

func DepositModal(ctx context.Context) {
	// Create a hotkey listener for the 'J' key
	hook.Register(hook.KeyDown, []string{"y"}, func(e hook.Event) {
		fmt.Println("Depositar pressionado")

		script := `
		window.dispatchEvent(new Event('openModalDeposit'))
		`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil), // Execute the script in the browser context
		)
		if errChrome != nil {
			fmt.Println("Erro ao abrir depósito:", errChrome)
		}

	})
}

func GotoHome(ctx context.Context) {
	hook.Register(hook.KeyDown, []string{"g"}, func(e hook.Event) {
		fmt.Println("Ir para lista de jogos pressionado")

		script := `
				window.dispatchEvent(new Event('gotoHome'))
			`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil), // Execute the script in the browser context
		)
		if errChrome != nil {
			fmt.Println("Erro ao tentar ir para a home:", errChrome)
		}

	})
}

// func ResetSession(ctx context.Context) {
// 	hook.Register(hook.KeyDown, []string{"r"}, func(e hook.Event) {
// 		fmt.Println("Saque solicitado - resetando a sessão")

// 		script := `
// 				window.dispatchEvent(new Event('clearAmount'))
// 			`
// 		errChrome := chromedp.Run(ctx,
// 			chromedp.Evaluate(script, nil), // Execute the script in the browser context
// 		)
// 		if errChrome != nil {
// 			fmt.Println("Erro ao tentar solicitar o saque:", errChrome)
// 		}

// 	})
// }

func ModalSaque(ctx context.Context) {
	hook.Register(hook.KeyDown, []string{"j"}, func(e hook.Event) {
		fmt.Println("abrir modal de saque")

		script := `
				window.dispatchEvent(new Event('openModalWithdrawall'))
			`
		errChrome := chromedp.Run(ctx,
			chromedp.Evaluate(script, nil), // Execute the script in the browser context
		)
		if errChrome != nil {
			fmt.Println("Erro ao tentar abrir o modal do saque:", errChrome)
		}

	})
}

func CloseProgram(ctx context.Context, cancel context.CancelFunc) {
	hook.Register(hook.KeyDown, []string{"f"}, func(e hook.Event) {
		fmt.Println("Fechando terminal")

		cancel()
		os.Exit(1)

	})
}
