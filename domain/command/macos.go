//go:build darwin
// +build darwin

package command

import "fmt"

type MacOSCommand struct {
}

func NewCommand() *MacOSCommand {
	return &MacOSCommand{}
}

func (c *MacOSCommand) DisableAltF4() error {
	fmt.Println("Não desabilita Alt+F4 no mac")

	return nil
}

func (c *MacOSCommand) EnableAltF4() error {
	fmt.Println("Não habilita Alt+F4 no mac")

	return nil
}

func (c *MacOSCommand) ChangeUser(user string, password string) error {
	fmt.Println("Não muda usuario no mac")

	return nil
}

func (c *MacOSCommand) SetNumLock(enabled bool) error {
	if enabled {
		fmt.Println("Mac does not support numlock, should enabled.")
	} else {
		fmt.Println("Mac does not support numlock, should disabled.")
	}

	return nil
}
