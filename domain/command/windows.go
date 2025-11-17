//go:build windows
// +build windows

package command

import (
	"fmt"
)

type WindowsCommand struct {
}

func NewCommand() *WindowsCommand {
	return &WindowsCommand{}
}

func (c *WindowsCommand) DisableAltF4() error {
	fmt.Println("Alt+F4 disabling is not supported on Windows.")
	// You could implement logic for interacting with system hooks if needed
	return nil
}

func (c *WindowsCommand) EnableAltF4() error {
	fmt.Println("Alt+F4 enabling is not supported on Windows.")
	// Implementation for enabling could go here
	return nil
}

func (c *WindowsCommand) ChangeUser(user string, password string) error {
	fmt.Println("User change is not directly supported on Windows.")
	// Implement Windows-specific logic for switching users if required
	return nil
}

func (c *WindowsCommand) SetNumLock(enabled bool) error {
	fmt.Println("SetNumLock its not supported on windows")
	return nil
}
