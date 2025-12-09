package domain

import (
	"fmt"
	"strings"
	"terminal/domain/command"
)

type Terminal struct {
	Id         string
	Name       string  // Derivado de UserName
	Status     string  // Sempre "Ativo" se conectou
	Url        string  // Configurado via env ou padrão
	BaseURL    string
	UserName   string  // user_name do banco
	Nid        int     // nid do banco
	Amount     float64 // amount do banco (saldo)
	FacilityId int     // facility_id do banco
	Command    command.Command
	PublicIp   string `json:"ip"`
	City       string
	Region     string
	Country    string
	Org        string
	Location   string
	Pin        string  // Pode ser configurado via env
	Account    string
	Session    *string
}

func (t Terminal) Slug() string {
	return strings.ReplaceAll(t.Name, " ", "-")
}

func (t Terminal) DisableKeys() {
	err := t.Command.DisableAltF4()
	if err != nil {
		panic(fmt.Sprintf("Erro ao desabilitar teclado: %s", err))
	}
}

func (t Terminal) EnableKeys() {
	t.Command.EnableAltF4()
}

func (t Terminal) GetSession() string {
	if t.Session == nil || *t.Session == "" {
		return ""
	}

	parts := strings.Split(*t.Session, "_")
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
