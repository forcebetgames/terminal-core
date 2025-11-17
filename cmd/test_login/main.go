package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"terminal/domain"
)

func main() {
	// Configurações de teste - MODIFIQUE AQUI
	baseURL := "https://dashboard.forcebetgames.com"
	username := "pc1"
	pin := "123"

	// Construir URL de login
	loginURL := fmt.Sprintf("%s/go-login?user_name=%s&pin=%s", baseURL, username, pin)

	fmt.Println("====================================")
	fmt.Println("🧪 TESTE DE LOGIN AUTOMÁTICO")
	fmt.Println("====================================")
	fmt.Printf("URL Base: %s\n", baseURL)
	fmt.Printf("Usuário: %s\n", username)
	fmt.Printf("PIN: %s\n", pin)
	fmt.Printf("URL de Login: %s\n", loginURL)
	fmt.Println("====================================")
	fmt.Println("\n🚀 Abrindo navegador...\n")

	// Criar contexto global
	globalCTX := context.TODO()

	// Abrir navegador
	_, cancel := domain.OpenBrowser(loginURL, globalCTX)
	defer cancel()

	fmt.Println("✅ Navegador aberto!")
	fmt.Println("📌 Pressione Ctrl+C para fechar")

	// Aguardar sinal de encerramento
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan

	fmt.Println("\n👋 Fechando...")
}
