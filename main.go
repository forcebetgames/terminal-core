package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"terminal/domain"
	"terminal/domain/command"
	"terminal/domain/keyboard"
	"terminal/logger"

	"github.com/denisbrodbeck/machineid"
	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

//go:embed .env
var embeddedEnv string

var safe_url = "https://google.com"

// ScreenInfo armazena informações sobre a tela
type ScreenInfo struct {
	Orientation string
	Width       int
	Height      int
}

// detectDisplayServer detecta se está rodando X11 ou Wayland
func detectDisplayServer() string {
	// 1. Verifica XDG_SESSION_TYPE (método mais confiável)
	if sessionType := os.Getenv("XDG_SESSION_TYPE"); sessionType != "" {
		return strings.ToLower(sessionType)
	}

	// 2. Verifica WAYLAND_DISPLAY
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return "wayland"
	}

	// 3. Verifica DISPLAY (X11)
	if os.Getenv("DISPLAY") != "" {
		return "x11"
	}

	// 4. Fallback: assume X11
	return "x11"
}

// detectScreenResolution detecta automaticamente a resolução e orientação da tela
func detectScreenResolution() ScreenInfo {
	info := ScreenInfo{
		Orientation: "portrait",
		Width:       1080,
		Height:      1920,
	}

	// 1. Verifica argumento de linha de comando para orientação
	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		if arg == "horizontal" || arg == "landscape" || arg == "-h" || arg == "--horizontal" {
			info.Orientation = "landscape"
		}
		if arg == "vertical" || arg == "portrait" || arg == "-v" || arg == "--vertical" {
			info.Orientation = "portrait"
		}
	}

	// 2. Verifica variável de ambiente SCREEN_ORIENTATION
	if orientation := os.Getenv("SCREEN_ORIENTATION"); orientation != "" {
		info.Orientation = orientation
	}

	// 3. Verifica se resolução foi definida manualmente via variáveis de ambiente
	if widthStr := os.Getenv("SCREEN_WIDTH"); widthStr != "" {
		if width, err := strconv.Atoi(widthStr); err == nil {
			info.Width = width
		}
	}
	if heightStr := os.Getenv("SCREEN_HEIGHT"); heightStr != "" {
		if height, err := strconv.Atoi(heightStr); err == nil {
			info.Height = height
		}
	}

	// Se resolução foi definida manualmente, retorna
	if os.Getenv("SCREEN_WIDTH") != "" && os.Getenv("SCREEN_HEIGHT") != "" {
		logger.Printf("🔧 Resolução manual: %dx%d\n", info.Width, info.Height)
		return info
	}

	// 4. Detecta display server
	displayServer := detectDisplayServer()
	logger.Printf("🖥️  Display Server: %s\n", strings.ToUpper(displayServer))

	// Estrutura para armazenar informações dos monitores
	type Monitor struct {
		name   string
		width  int
		height int
		area   int
	}

	var monitors []Monitor

	// 5. Tenta detectar resolução baseado no display server
	if displayServer == "wayland" {
		logger.Println("🌊 Detectando resolução no Wayland...")

		// Método 1: wlr-randr (para wlroots compositors: Sway, Hyprland, etc)
		cmd := exec.Command("wlr-randr")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			var currentMonitor string

			for _, line := range lines {
				line = strings.TrimSpace(line)

				// Detecta nome do monitor
				if len(line) > 0 && !strings.HasPrefix(line, " ") && strings.Contains(line, " ") {
					parts := strings.Fields(line)
					if len(parts) > 0 {
						currentMonitor = parts[0]
					}
				}

				// Detecta resolução atual (linha com asterisco *)
				if strings.Contains(line, "current") || strings.Contains(line, "*") {
					parts := strings.Fields(line)
					for _, part := range parts {
						if strings.Contains(part, "x") {
							dims := strings.Split(strings.TrimSuffix(part, ","), "x")
							if len(dims) == 2 {
								width, errW := strconv.Atoi(dims[0])
								height, errH := strconv.Atoi(dims[1])
								if errW == nil && errH == nil && width > 0 && height > 0 {
									monitors = append(monitors, Monitor{
										name:   currentMonitor,
										width:  width,
										height: height,
										area:   width * height,
									})
									logger.Printf("   ✓ Monitor %s: %dx%d (wlr-randr)\n", currentMonitor, width, height)
									break
								}
							}
						}
					}
				}
			}
		}

		// Método 2: gnome-randr (para GNOME/Mutter no Wayland)
		if len(monitors) == 0 {
			cmd := exec.Command("gnome-randr")
			output, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					if strings.Contains(line, "*") || strings.Contains(line, "current") {
						parts := strings.Fields(line)
						for _, part := range parts {
							if strings.Contains(part, "x") {
								dims := strings.Split(part, "x")
								if len(dims) == 2 {
									width, errW := strconv.Atoi(dims[0])
									height, errH := strconv.Atoi(dims[1])
									if errW == nil && errH == nil && width > 0 && height > 0 {
										monitors = append(monitors, Monitor{
											name:   "GNOME",
											width:  width,
											height: height,
											area:   width * height,
										})
										logger.Printf("   ✓ Monitor: %dx%d (gnome-randr)\n", width, height)
										break
									}
								}
							}
						}
					}
				}
			}
		}

		// Método 3: kscreen-doctor (para KDE Plasma no Wayland)
		if len(monitors) == 0 {
			cmd := exec.Command("kscreen-doctor", "-o")
			output, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(output), "\n")
				for _, line := range lines {
					if strings.Contains(line, "Output:") && strings.Contains(line, "enabled") {
						// Procura por resolução na linha seguinte ou mesma linha
						parts := strings.Fields(line)
						for _, part := range parts {
							if strings.Contains(part, "x") && strings.Contains(part, "@") {
								resolution := strings.Split(part, "@")[0]
								dims := strings.Split(resolution, "x")
								if len(dims) == 2 {
									width, errW := strconv.Atoi(dims[0])
									height, errH := strconv.Atoi(dims[1])
									if errW == nil && errH == nil && width > 0 && height > 0 {
										monitors = append(monitors, Monitor{
											name:   "KDE",
											width:  width,
											height: height,
											area:   width * height,
										})
										logger.Printf("   ✓ Monitor: %dx%d (kscreen-doctor)\n", width, height)
										break
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// 6. Tenta xrandr (funciona em X11 e alguns Wayland com XWayland)
	if len(monitors) == 0 {
		if displayServer == "x11" {
			logger.Println("🪟 Detectando resolução no X11...")
		} else {
			logger.Println("   Tentando xrandr via XWayland...")
		}

		cmd := exec.Command("xrandr", "--current")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")

			for _, line := range lines {
				if strings.Contains(line, " connected") {
					parts := strings.Fields(line)
					monitorName := parts[0]

					for _, part := range parts {
						if strings.Contains(part, "x") && (strings.Contains(part, "+") || len(strings.Split(part, "x")) == 2) {
							resolution := strings.Split(part, "+")[0]
							dims := strings.Split(resolution, "x")
							if len(dims) == 2 {
								width, errW := strconv.Atoi(dims[0])
								height, errH := strconv.Atoi(dims[1])
								if errW == nil && errH == nil && width > 0 && height > 0 {
									monitors = append(monitors, Monitor{
										name:   monitorName,
										width:  width,
										height: height,
										area:   width * height,
									})
									logger.Printf("   ✓ Monitor %s: %dx%d (xrandr)\n", monitorName, width, height)
									break
								}
							}
						}
					}
				}
			}
		}
	}

	// 7. Se encontrou monitores, seleciona o de MAIOR resolução (área)
	if len(monitors) > 0 {
		largest := monitors[0]
		for _, m := range monitors {
			if m.area > largest.area {
				largest = m
			}
		}

		if len(monitors) > 1 {
			logger.Printf("📺 Múltiplos monitores detectados (%d)\n", len(monitors))
			logger.Printf("✅ Usando monitor: %s (%dx%d - maior resolução)\n", largest.name, largest.width, largest.height)
		} else {
			logger.Printf("✅ Resolução detectada: %dx%d\n", largest.width, largest.height)
		}

		// Define resolução detectada
		info.Width = largest.width
		info.Height = largest.height

		// Detecta orientação baseada na resolução (se não foi especificada manualmente)
		if os.Getenv("SCREEN_ORIENTATION") == "" && len(os.Args) <= 1 {
			if largest.width > largest.height {
				info.Orientation = "landscape"
				logger.Printf("📐 Orientação: HORIZONTAL (landscape)\n")
			} else {
				info.Orientation = "portrait"
				logger.Printf("📐 Orientação: VERTICAL (portrait)\n")
			}
		}

		return info
	}

	// 8. Fallback: Tenta detectar via /sys/class/graphics (framebuffer - funciona em ambos)
	logger.Println("   Tentando detectar via framebuffer...")
	fbFiles := []string{
		"/sys/class/graphics/fb0/virtual_size",
		"/sys/class/drm/card0/card0-HDMI-A-1/modes",
		"/sys/class/drm/card0/card0-DP-1/modes",
	}

	for _, fbFile := range fbFiles {
		fbOutput, fbErr := exec.Command("cat", fbFile).Output()
		if fbErr == nil {
			content := strings.TrimSpace(string(fbOutput))

			// virtual_size format: "width,height"
			if strings.Contains(content, ",") {
				parts := strings.Split(content, ",")
				if len(parts) == 2 {
					width, _ := strconv.Atoi(parts[0])
					height, _ := strconv.Atoi(parts[1])
					if width > 0 && height > 0 {
						info.Width = width
						info.Height = height
						logger.Printf("✅ Resolução detectada (framebuffer): %dx%d\n", width, height)

						// Detecta orientação (se não foi especificada manualmente)
						if os.Getenv("SCREEN_ORIENTATION") == "" && len(os.Args) <= 1 {
							if width > height {
								info.Orientation = "landscape"
								logger.Printf("📐 Orientação: HORIZONTAL (landscape)\n")
							} else {
								info.Orientation = "portrait"
								logger.Printf("📐 Orientação: VERTICAL (portrait)\n")
							}
						}

						return info
					}
				}
			}

			// modes format: "widthxheight"
			if strings.Contains(content, "x") {
				lines := strings.Split(content, "\n")
				if len(lines) > 0 {
					dims := strings.Split(lines[0], "x")
					if len(dims) == 2 {
						width, _ := strconv.Atoi(dims[0])
						height, _ := strconv.Atoi(dims[1])
						if width > 0 && height > 0 {
							info.Width = width
							info.Height = height
							logger.Printf("✅ Resolução detectada (DRM): %dx%d\n", width, height)

							// Detecta orientação
							if os.Getenv("SCREEN_ORIENTATION") == "" && len(os.Args) <= 1 {
								if width > height {
									info.Orientation = "landscape"
									logger.Printf("📐 Orientação: HORIZONTAL (landscape)\n")
								} else {
									info.Orientation = "portrait"
									logger.Printf("📐 Orientação: VERTICAL (portrait)\n")
								}
							}

							return info
						}
					}
				}
			}
		}
	}

	// 9. Fallback final: valores padrão baseados na orientação
	logger.Println("⚠️  Não foi possível detectar resolução, usando valores padrão")
	if info.Orientation == "landscape" {
		info.Width = 1920
		info.Height = 1080
		logger.Printf("📐 Usando padrão: %dx%d (landscape)\n", info.Width, info.Height)
	} else {
		info.Width = 1080
		info.Height = 1920
		logger.Printf("📐 Usando padrão: %dx%d (portrait)\n", info.Width, info.Height)
	}

	return info
}

func main() {
	// Parse command line flags
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Initialize logger
	logger.Init(*debugFlag)

	globalCTX := context.TODO()

	logger.Println("====================================")
	logger.Println("🎮 TERMINAL DE JOGOS - FORCEBET")
	logger.Println("====================================")
	logger.Println("")
	logger.Println("Iniciando o sistema...")

	// Carrega .env PRIMEIRO
	envMap, errEnv := godotenv.Unmarshal(embeddedEnv)
	if errEnv != nil {
		log.Printf("Error loading .env file: %s", errEnv)
		fmt.Scanln()
		os.Exit(1)
	}

	// Set the variables in the runtime environment
	for key, value := range envMap {
		os.Setenv(key, value)
	}

	// DEPOIS detecta resolução e orientação automaticamente
	screenInfo := detectScreenResolution()
	os.Setenv("SCREEN_ORIENTATION", screenInfo.Orientation)
	os.Setenv("SCREEN_WIDTH", strconv.Itoa(screenInfo.Width))
	os.Setenv("SCREEN_HEIGHT", strconv.Itoa(screenInfo.Height))

	if screenInfo.Orientation == "landscape" {
		logger.Printf("📐 Modo: HORIZONTAL (Landscape) - %dx%d\n", screenInfo.Width, screenInfo.Height)
		logger.Println("🎯 Jogo: Empire 🏛️")
	} else {
		logger.Printf("📐 Modo: VERTICAL (Portrait) - %dx%d\n", screenInfo.Width, screenInfo.Height)
		logger.Println("🎯 Jogo: Tigrinho 🐯")
	}
	logger.Println("")
	logger.Println("====================================")

	///////////////////////////////////////////////////////
	trmEnvType := "dev"

	if os.Getenv("MYSQL_DATABASE") == "defaultdb" {
		trmEnvType = "prod"
	}
	///////////////////////////////////////////////////////

	logger.Printf("Se conectando ao banco de dados de %s \n", trmEnvType)
	db := domain.NewDatabaseConnection()

	logger.Println("Pegando a numeração do terminal...")
	machineId, machineIdErr := machineid.ID()
	if machineIdErr != nil {
		log.Printf("Não foi possível identificar a numeração de serie do terminal")
		fmt.Scanln()
		os.Exit(1)
	}
	logger.Println("====================================")
	logger.Printf("🔍 DEBUG - Machine ID detectado: %s\n", machineId)
	logger.Println("====================================")

	// Query com TODOS os campos necessários
	query := "SELECT id, user_name, pin, name, amount, status, IFNULL(url, ?), nid, facility_id, session_id FROM trm_terminal WHERE id = ?"
	logger.Printf("📊 Executando query:\n%s\n", query)
	logger.Printf("📊 Parâmetros: safe_url='%s', machineId='%s'\n", safe_url, machineId)

	row := db.QueryRow(query, safe_url, machineId)

	var terminal domain.Terminal
	var amountCents int // amount está em centavos no banco

	if err := row.Scan(&terminal.Id, &terminal.UserName, &terminal.Pin, &terminal.Name, &amountCents, &terminal.Status, &terminal.Url, &terminal.Nid, &terminal.FacilityId, &terminal.Session); err != nil {
		logger.Println("====================================")
		logger.Println("❌ ERRO ao fazer Scan dos dados do terminal!")
		logger.Printf("Machine ID buscado: %s\n", machineId)
		logger.Printf("Erro detalhado: %v\n", err)
		logger.Println("====================================")

		// Tenta verificar se o registro existe
		var count int
		countQuery := "SELECT COUNT(*) FROM trm_terminal WHERE id = ?"
		if err := db.QueryRow(countQuery, machineId).Scan(&count); err == nil {
			logger.Printf("✓ Terminal EXISTE no banco (count: %d)\n", count)
			logger.Println("  → Verifique os tipos de dados dos campos!")
		} else {
			logger.Printf("✗ Terminal NÃO existe no banco\n")
		}

		logger.Println("====================================")
		fmt.Scanln()
		os.Exit(1)
	}

	// Converte amount de centavos para reais (R$ 1.150,84 = 115084 centavos)
	terminal.Amount = float64(amountCents) / 100.0

	logger.Println("====================================")
	logger.Println("📊 DEBUG - Dados do Terminal:")
	logger.Printf("  ID (Machine ID): %s\n", terminal.Id)
	logger.Printf("  Nome: %s\n", terminal.Name)
	logger.Printf("  Username: %s\n", terminal.UserName)
	logger.Printf("  NID: %d\n", terminal.Nid)
	logger.Printf("  Saldo: R$ %.2f\n", terminal.Amount)
	logger.Printf("  Facility ID: %d\n", terminal.FacilityId)
	logger.Printf("  Status: %s\n", terminal.Status)
	logger.Printf("  URL: %s\n", terminal.Url)
	logger.Printf("  PIN: %s\n", terminal.Pin)
	if terminal.Session != nil {
		logger.Printf("  Session ID: %s\n", *terminal.Session)
	} else {
		logger.Println("  Session ID: NULL")
	}
	logger.Println("====================================")

	logger.Println("Validando o terminal")
	if terminal.Status != "Ativo" {
		log.Printf("Terminal inativo")
		fmt.Scanln()
		os.Exit(1)
	}

	publicIpErr := domain.GetIPInfo(&terminal)
	// if publicIpErr != nil {
	// 	// log.Fatalf("%s", fmt.Sprintf("Erro ao identificar IP publico, %s", publicIpErr))
	// 	fmt.Scanln()
	// 	// os.Exit(1)
	// }

	if publicIpErr == nil {
		_, errUpdate := db.Exec("UPDATE trm_terminal set public_ip = ? WHERE id = ? ", terminal.PublicIp, terminal.Id)
		if errUpdate != nil {
			log.Printf("%s", fmt.Sprintf("Erro ao atualizar o IP publico, %s", errUpdate))
			fmt.Scanln()
			os.Exit(1)
		}
	}

	terminal.Command = command.NewCommand()
	terminal.DisableKeys()

	logger.Println("Estabelecendo comunicação")

	if terminal.Url != safe_url {
		terminal.BaseURL = terminal.Url
		terminal.Url += fmt.Sprintf("/go-login?user_name=%s&pin=%s", terminal.UserName, terminal.Pin)
	}

	logger.Println("====================================")
	logger.Println("🌐 DEBUG - URLs Construídas:")
	logger.Printf("  URL Base: %s\n", terminal.BaseURL)
	logger.Printf("  URL Final (com login): %s\n", terminal.Url)
	logger.Println("====================================")
	logger.Println("🚀 Abrindo navegador Chromium...")
	logger.Println("====================================")

	ctx, cancelBrowser := domain.OpenBrowser(terminal.Url, globalCTX)

	logger.Println("====================================")
	logger.Println("✅ Navegador aberto com sucesso!")
	logger.Println("====================================")
	logger.Println("🎹 Inicializando sistema de input...")

	inputHandler, err := keyboard.NewInputHandler()
	if err != nil {
		log.Printf("❌ Erro ao inicializar input handler: %v", err)
		logger.Println("💡 Dica: Se estiver usando Wayland, execute com sudo!")
		fmt.Scanln()
		os.Exit(1)
	}

	logger.Println("====================================")
	logger.Println("📌 Registrando atalhos de teclado...")
	logger.Println("====================================")

	command.DepositModal(ctx, inputHandler)
	command.GotoHome(ctx, inputHandler)
	command.ModalSaque(ctx, inputHandler)
	command.ConfirmWithdrawal(ctx, inputHandler)
	command.ShowReport(ctx, inputHandler)
	command.CloseProgram(ctx, cancelBrowser, inputHandler)
	command.IncreaseBet(ctx, inputHandler)
	command.DecreaseBet(ctx, inputHandler)
	command.ChangeGame(ctx, inputHandler)
	command.StartGame(ctx, inputHandler)

	if terminal.GetSession() == "PIX" {
		terminal.Command.SetNumLock(true)
		logger.Println("🔢 NumLock: HABILITADO (sessão PIX)")
	} else {
		terminal.Command.SetNumLock(false)
		logger.Println("🔢 NumLock: DESABILITADO (sessão CASH)")
	}

	logger.Println("====================================")
	logger.Println("⌨️  Atalhos de teclado disponíveis:")
	logger.Println("  Y - PIX / Abrir modal de depósito")
	logger.Println("  J - SACAR / Abrir modal de saque")
	logger.Println("  R - CONFIRMAR SAQUE")
	logger.Println("  B - RELATÓRIO de movimentações e cobranças")
	logger.Println("  G - Ir para home/jogos")
	logger.Println("  F - Fechar programa")
	logger.Println("")
	logger.Println("🎮 CONTROLES DO JOGO:")
	logger.Println("  + (=) - Aumentar aposta")
	logger.Println("  -     - Diminuir aposta")
	logger.Println("  G ou ↑ ou O - Mudar jogo")
	logger.Println("  SPACE - Iniciar/Girar jogo")
	logger.Println("")
	logger.Println("💰 TECLAS DE DINHEIRO (aceita AMBAS - teclas normais OU numpad):")
	logger.Println("  2 ou Numpad 2  → R$ 2")
	logger.Println("  3 ou Numpad 3  → R$ 5")
	logger.Println("  4 ou Numpad 4  → R$ 10")
	logger.Println("  5 ou Numpad 5  → R$ 20")
	logger.Println("  6 ou Numpad 6  → R$ 50")
	logger.Println("  7 ou Numpad 7  → R$ 100")
	logger.Println("====================================")
	logger.Println("🔌 Conectando ao Pusher (websocket)...")

	pusher := domain.NewPusher()
	depositDoneErr := pusher.Listen(&terminal, "deposit_done", func(data map[string]interface{}) {
		session, sessionOK := data["session"]
		if !sessionOK {
			return
		}

		if session == "PIX" {
			terminal.Command.SetNumLock(true)
		}

		if session == "CASH" {
			terminal.Command.SetNumLock(false)
		}
	})

	if depositDoneErr != nil {
		log.Printf("%s", fmt.Sprintf("Erro ao ler socket deposito, %s", depositDoneErr))
		fmt.Scanln()
		os.Exit(1)
	}

	amountErr := pusher.Listen(&terminal, "amount", func(data map[string]interface{}) {
		transactionType, transactionTypeOk := data["transactionType"]
		if !transactionTypeOk {
			return
		}

		if transactionType != "inactivity" {
			return
		}

		terminal.Command.SetNumLock(false)
	})

	if amountErr != nil {
		log.Printf("%s", fmt.Sprintf("Erro ao ler socket saldo, %s", amountErr))
		fmt.Scanln()
		os.Exit(1)
	}

	logger.Println("✅ Pusher conectado!")
	logger.Println("====================================")

	payment := domain.NewPaymentCash(func(brlCount int) {
		logger.Println("====================================")
		logger.Printf("🚀 CALLBACK EXECUTADO - Valor: R$ %d\n", brlCount)
		logger.Println("====================================")

		data := map[string]interface{}{
			"amount":     brlCount,
			"notas":      map[string]int{strconv.Itoa(brlCount): 1},
			"terminalId": terminal.Id,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Printf("❌ #1 - Erro ao marshalar JSON: %v\n", err)
			return
		}

		postURL := fmt.Sprintf("%s/api/hooks/pnr/deposit_cash", terminal.BaseURL)
		logger.Printf("📡 Enviando POST para: %s\n", postURL)
		logger.Printf("📦 Payload: %s\n", string(jsonData))

		req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Printf("❌ #2 - Erro ao criar requisição: %v\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Printf("❌ #3 - Erro ao executar requisição: %v\n", err)
			return
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Printf("❌ #B - Erro ao ler resposta: %v\n", err)
			return
		}

		logger.Printf("📥 Status HTTP: %d\n", resp.StatusCode)
		logger.Printf("📥 Resposta do servidor: %s\n", string(respBody))

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			logger.Printf("❌ #4 - Resposta de erro do servidor (status %d)\n", resp.StatusCode)
			logger.Printf("   Corpo: %s\n", string(respBody))
		} else {
			logger.Println("✅ Pagamento enviado com sucesso!")
		}

		logger.Println("====================================")

	}, inputHandler)
	payment.Start()

	// Start listening for keyboard events
	inputHandler.StartListening()

	logger.Println("====================================")
	logger.Println("")
	logger.Println("🎮 SISTEMA PRONTO!")
	logger.Println("Pressione Ctrl+C para encerrar")
	logger.Println("")
	logger.Println("💡 Tente pressionar teclas 2-7 ou Numpad 2-7")
	logger.Println("   Você verá logs detalhados no console!")
	logger.Println("")
	logger.Println("====================================")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan

}
