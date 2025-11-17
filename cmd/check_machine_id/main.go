package main

import (
	"fmt"
	"log"
	"os"

	"github.com/denisbrodbeck/machineid"
)

func main() {
	fmt.Println("====================================")
	fmt.Println("🔍 VERIFICADOR DE MACHINE ID")
	fmt.Println("====================================")
	fmt.Println("")

	machineId, err := machineid.ID()
	if err != nil {
		log.Printf("Erro ao obter Machine ID: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Machine ID desta máquina: %s\n", machineId)
	fmt.Println("")
	fmt.Println("Use este ID para buscar/criar o terminal no banco de dados.")
	fmt.Println("")
	fmt.Println("Query SQL para verificar:")
	fmt.Printf("  SELECT * FROM trm_terminal WHERE id = '%s';\n", machineId)
	fmt.Println("")
	fmt.Println("====================================")
}
