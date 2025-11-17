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
	"os/signal"
	"strconv"
	"syscall"
	"terminal/domain"
	"terminal/domain/command"
	"terminal/domain/keyboard"

	"github.com/denisbrodbeck/machineid"
	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
	hook "github.com/robotn/gohook"
)

//go:embed .env
var embeddedEnv string

var safe_url = "https://google.com"

func init() {
	// Start the global hook listener
	go func() {
		fmt.Println("Initializing hook system...")
		s := hook.Start()
		<-hook.Process(s) // Keep the hook system running
	}()
}

func main() {

	globalCTX := context.TODO()

	fmt.Println("Iniciando o sistema...")
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
	fmt.Println("📌 Registrando atalhos de teclado...")
	fmt.Println("====================================")

	command.DepositModal(ctx)
	command.GotoHome(ctx)
	// command.ResetSession(ctx)

	command.ModalSaque(ctx)
	command.CloseProgram(ctx, cancelBrowser)

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

	inputHandler := keyboard.NewRealInputHandler()

	fmt.Println("✅ Pusher conectado!")
	fmt.Println("====================================")

	payment := domain.NewPaymentCash(func(brlCount int) {
		data := map[string]interface{}{
			"amount":     brlCount,
			"notas":      map[string]int{strconv.Itoa(brlCount): 1},
			"terminalId": terminal.Id,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("#1 - Erro ao enviar a nota ao sistema")
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/hooks/pnr/deposit_cash", terminal.BaseURL), bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("#2 - Erro ao enviar a nota ao sistema")
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("#3 - Erro ao enviar a nota ao sistema")
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("#B - Erro ao enviar a nota ao sistema")
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			fmt.Println("#4 - Erro ao enviar a nota ao sistema", respBody)
		}

		defer resp.Body.Close()

	}, inputHandler)
	payment.Start()

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
