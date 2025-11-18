package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
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

	"github.com/denisbrodbeck/machineid"
	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

//go:embed .env
var embeddedEnv string

var safe_url = "https://google.com"

// detectScreenOrientation detecta automaticamente a orientação da tela
func detectScreenOrientation() string {
	// 1. Verifica argumento de linha de comando
	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		if arg == "horizontal" || arg == "landscape" || arg == "-h" || arg == "--horizontal" {
			return "landscape"
		}
		if arg == "vertical" || arg == "portrait" || arg == "-v" || arg == "--vertical" {
			return "portrait"
		}
	}

	// 2. Verifica variável de ambiente SCREEN_ORIENTATION
	if orientation := os.Getenv("SCREEN_ORIENTATION"); orientation != "" {
		return orientation
	}

	// 3. Tenta detectar automaticamente usando xrandr (X11/Wayland)
	cmd := exec.Command("xrandr", "--current")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")

		// Estrutura para armazenar informações dos monitores
		type Monitor struct {
			name   string
			width  int
			height int
			area   int
		}

		var monitors []Monitor

		// Coleta informações de todos os monitores conectados
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
								break
							}
						}
					}
				}
			}
		}

		// Se encontrou monitores, seleciona o de MAIOR resolução (área)
		if len(monitors) > 0 {
			largest := monitors[0]
			for _, m := range monitors {
				if m.area > largest.area {
					largest = m
				}
			}

			if len(monitors) > 1 {
				fmt.Printf("🖥️  Múltiplos monitores detectados (%d)\n", len(monitors))
				fmt.Printf("📺 Usando monitor: %s (%dx%d - maior resolução)\n", largest.name, largest.width, largest.height)
			}

			if largest.width > largest.height {
				fmt.Printf("🔍 Resolução detectada: %dx%d → HORIZONTAL\n", largest.width, largest.height)
				return "landscape"
			} else {
				fmt.Printf("🔍 Resolução detectada: %dx%d → VERTICAL\n", largest.width, largest.height)
				return "portrait"
			}
		}
	}

	// 4. Fallback: Tenta detectar via /sys/class/graphics (Wayland/console)
	fbOutput, fbErr := exec.Command("cat", "/sys/class/graphics/fb0/virtual_size").Output()
	if fbErr == nil {
		parts := strings.Split(strings.TrimSpace(string(fbOutput)), ",")
		if len(parts) == 2 {
			width, _ := strconv.Atoi(parts[0])
			height, _ := strconv.Atoi(parts[1])
			if width > height {
				fmt.Printf("🔍 Resolução detectada (fb0): %dx%d → HORIZONTAL\n", width, height)
				return "landscape"
			} else {
				fmt.Printf("🔍 Resolução detectada (fb0): %dx%d → VERTICAL\n", width, height)
				return "portrait"
			}
		}
	}

	// 5. Fallback final: portrait (vertical)
	fmt.Println("⚠️  Não foi possível detectar orientação, usando padrão: VERTICAL")
	return "portrait"
}

func main() {

	globalCTX := context.TODO()

	fmt.Println("====================================")
	fmt.Println("🎮 TERMINAL DE JOGOS - FORCEBET")
	fmt.Println("====================================")
	fmt.Println("")
	fmt.Println("Iniciando o sistema...")

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

	// DEPOIS detecta orientação (sobrescreve .env se necessário)
	orientation := detectScreenOrientation()
	os.Setenv("SCREEN_ORIENTATION", orientation)

	if orientation == "landscape" {
		fmt.Println("📐 Modo: HORIZONTAL (Landscape) - 1920x1080")
		fmt.Println("🎯 Jogo: Empire 🏛️")
	} else {
		fmt.Println("📐 Modo: VERTICAL (Portrait) - 1080x1920")
		fmt.Println("🎯 Jogo: Tigrinho 🐯")
	}
	fmt.Println("")
	fmt.Println("====================================")

	///////////////////////////////////////////////////////
	trmEnvType := "dev"

	if os.Getenv("MYSQL_DATABASE") == "terminal_prod" {
		trmEnvType = "prod"
	}
	///////////////////////////////////////////////////////

	fmt.Printf("Se conectando ao banco de dados de %s \n", trmEnvType)
	db := domain.NewDatabaseConnection()

	fmt.Println("Pegando a numeração do terminal...")
	machineId, machineIdErr := machineid.ID()
	if machineIdErr != nil {
		log.Printf("Não foi possível identificar a numeração de serie do terminal")
		fmt.Scanln()
		os.Exit(1)
	}
	fmt.Println("====================================")
	fmt.Printf("🔍 DEBUG - Machine ID detectado: %s\n", machineId)
	fmt.Println("====================================")

	row := db.QueryRow("SELECT id, name, status, IFNULL(url,?), user_name,pin,session_id from trm_terminal where id = ?", safe_url, machineId)

	var terminal domain.Terminal
	if err := row.Scan(&terminal.Id, &terminal.Name, &terminal.Status, &terminal.Url, &terminal.UserName, &terminal.Pin, &terminal.Session); err != nil {
		log.Printf("%s", fmt.Sprintf("Terminal não encontrado: %s", machineId))
		fmt.Println("====================================")
		fmt.Println("❌ ERRO: Terminal não existe no banco de dados!")
		fmt.Printf("Machine ID buscado: %s\n", machineId)
		fmt.Println("====================================")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Println("====================================")
	fmt.Println("📊 DEBUG - Dados do Terminal:")
	fmt.Printf("  ID: %s\n", terminal.Id)
	fmt.Printf("  Nome: %s\n", terminal.Name)
	fmt.Printf("  Status: %s\n", terminal.Status)
	fmt.Printf("  URL do banco: %s\n", terminal.Url)
	fmt.Printf("  Username: %s\n", terminal.UserName)
	fmt.Printf("  PIN: %s\n", terminal.Pin)
	if terminal.Session != nil {
		fmt.Printf("  Session ID: %s\n", *terminal.Session)
	} else {
		fmt.Println("  Session ID: NULL")
	}
	fmt.Println("====================================")

	fmt.Println("Validando o terminal")
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

	fmt.Println("Estabelecendo comunicação")

	if terminal.Url != safe_url {
		terminal.BaseURL = terminal.Url
		terminal.Url += fmt.Sprintf("/go-login?user_name=%s&pin=%s", terminal.UserName, terminal.Pin)
	}

	fmt.Println("====================================")
	fmt.Println("🌐 DEBUG - URLs Construídas:")
	fmt.Printf("  URL Base: %s\n", terminal.BaseURL)
	fmt.Printf("  URL Final (com login): %s\n", terminal.Url)
	fmt.Println("====================================")
	fmt.Println("🚀 Abrindo navegador Chromium...")
	fmt.Println("====================================")

	ctx, cancelBrowser := domain.OpenBrowser(terminal.Url, globalCTX)

	fmt.Println("====================================")
	fmt.Println("✅ Navegador aberto com sucesso!")
	fmt.Println("====================================")
	fmt.Println("🎹 Inicializando sistema de input...")

	inputHandler, err := keyboard.NewInputHandler()
	if err != nil {
		log.Printf("❌ Erro ao inicializar input handler: %v", err)
		fmt.Println("💡 Dica: Se estiver usando Wayland, execute com sudo!")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Println("====================================")
	fmt.Println("📌 Registrando atalhos de teclado...")
	fmt.Println("====================================")

	command.DepositModal(ctx, inputHandler)
	command.GotoHome(ctx, inputHandler)
	command.ModalSaque(ctx, inputHandler)
	command.CloseProgram(ctx, cancelBrowser, inputHandler)

	if terminal.GetSession() == "PIX" {
		terminal.Command.SetNumLock(true)
		fmt.Println("🔢 NumLock: HABILITADO (sessão PIX)")
	} else {
		terminal.Command.SetNumLock(false)
		fmt.Println("🔢 NumLock: DESABILITADO (sessão CASH)")
	}

	fmt.Println("====================================")
	fmt.Println("⌨️  Atalhos de teclado disponíveis:")
	fmt.Println("  Y - Abrir modal de depósito")
	fmt.Println("  G - Ir para home/jogos")
	fmt.Println("  J - Abrir modal de saque")
	fmt.Println("  F - Fechar programa")
	fmt.Println("")
	fmt.Println("💰 TECLAS DE DINHEIRO (aceita AMBAS):")
	fmt.Println("  2 ou Numpad 2  → R$ 2")
	fmt.Println("  3 ou Numpad 3  → R$ 5")
	fmt.Println("  4 ou Numpad 4  → R$ 10")
	fmt.Println("  5 ou Numpad 5  → R$ 20")
	fmt.Println("  6 ou Numpad 6  → R$ 50")
	fmt.Println("  7 ou Numpad 7  → R$ 100")
	fmt.Println("====================================")
	fmt.Println("🔌 Conectando ao Pusher (websocket)...")

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

	fmt.Println("✅ Pusher conectado!")
	fmt.Println("====================================")

	payment := domain.NewPaymentCash(func(brlCount int) {
		fmt.Println("====================================")
		fmt.Printf("🚀 CALLBACK EXECUTADO - Valor: R$ %d\n", brlCount)
		fmt.Println("====================================")

		data := map[string]interface{}{
			"amount":     brlCount,
			"notas":      map[string]int{strconv.Itoa(brlCount): 1},
			"terminalId": terminal.Id,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("❌ #1 - Erro ao marshalar JSON: %v\n", err)
			return
		}

		postURL := fmt.Sprintf("%s/api/hooks/pnr/deposit_cash", terminal.BaseURL)
		fmt.Printf("📡 Enviando POST para: %s\n", postURL)
		fmt.Printf("📦 Payload: %s\n", string(jsonData))

		req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("❌ #2 - Erro ao criar requisição: %v\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ #3 - Erro ao executar requisição: %v\n", err)
			return
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("❌ #B - Erro ao ler resposta: %v\n", err)
			return
		}

		fmt.Printf("📥 Status HTTP: %d\n", resp.StatusCode)
		fmt.Printf("📥 Resposta do servidor: %s\n", string(respBody))

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			fmt.Printf("❌ #4 - Resposta de erro do servidor (status %d)\n", resp.StatusCode)
			fmt.Printf("   Corpo: %s\n", string(respBody))
		} else {
			fmt.Println("✅ Pagamento enviado com sucesso!")
		}

		fmt.Println("====================================")

	}, inputHandler)
	payment.Start()

	// Start listening for keyboard events
	inputHandler.StartListening()

	fmt.Println("====================================")
	fmt.Println("")
	fmt.Println("🎮 SISTEMA PRONTO!")
	fmt.Println("Pressione Ctrl+C para encerrar")
	fmt.Println("")
	fmt.Println("💡 Tente pressionar teclas 2-7 ou Numpad 2-7")
	fmt.Println("   Você verá logs detalhados no console!")
	fmt.Println("")
	fmt.Println("====================================")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan

}
