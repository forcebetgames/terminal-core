package domain

import (
	"context"
	"log"
	"os"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func OpenBrowser(url string, globalCTX context.Context) (context.Context, context.CancelFunc) {
	headers := map[string]any{
		"X-Custom-keyevent": "potato",
	}

	// Detectar se DevTools deve estar habilitado
	enableDevTools := os.Getenv("ENABLE_DEVTOOLS") == "true"

	// Detectar display server para otimizações específicas
	displayServer := os.Getenv("XDG_SESSION_TYPE")
	isWayland := displayServer == "wayland"
	isX11 := displayServer == "x11" || os.Getenv("DISPLAY") != ""

	if enableDevTools {
		log.Println("🐛 DevTools HABILITADO - Modo Debug")
		log.Println("   Pressione F12 para abrir DevTools")
	} else {
		log.Println("🔒 DevTools DESABILITADO - Modo Produção (Kiosk)")
	}

	log.Println("🎮 Modo OTIMIZADO PARA JOGOS - Aceleração GPU Forçada")
	log.Println("   (Compatível com controles touch/teclado/mouse + Performance Máxima)")

	if isWayland {
		log.Println("🌊 Display: Wayland (otimizações completas)")
	} else if isX11 {
		log.Println("🪟 Display: X11/Xorg (otimizações compatíveis)")
	}

	// ✅ VERSÃO MINIMALISTA - Baseada na versão antiga fluida (11 flags apenas)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("kiosk", !enableDevTools),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("enable-automation", false),
	)
	log.Println("⚡ MODO MINIMALISTA: Configuração original fluida (11 flags)")

	// Se DevTools estiver habilitado, adiciona flags específicas
	if enableDevTools {
		opts = append(opts,
			// Abre DevTools automaticamente
			chromedp.Flag("auto-open-devtools-for-tabs", true),
			// Inicia o navegador maximizado (não fullscreen kiosk)
			chromedp.Flag("start-maximized", true),
		)
	}

	ctx, _ := chromedp.NewExecAllocator(globalCTX, opts...)

	ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)

	// Configura headers customizados
	if err := chromedp.Run(ctx, network.SetExtraHTTPHeaders(network.Headers(headers))); err != nil {
		log.Fatalf("Failed to set headers: %v", err)
		os.Exit(1)
	}

	// Navegação simples - como na versão antiga
	err := chromedp.Run(ctx, chromedp.Navigate(url))
	if err != nil {
		log.Fatalf("Error navigating to URL: %v", err)
		os.Exit(1)
	}

	log.Println("✅ Navegador iniciado (versão minimalista - igual à fluida)")

	return ctx, cancel
}
