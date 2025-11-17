//go:build linux
// +build linux

package command

import (
	"fmt"
	"os"
	"os/exec"
)

type UbuntuCommand struct {
}

func NewCommand() *UbuntuCommand {
	return &UbuntuCommand{}
}

func (c *UbuntuCommand) DisableAltF4() error {
	errAlt := exec.Command("xmodmap", "-e", "keycode 64 = NoSymbol").Run()
	if errAlt != nil {
		return errAlt
	}

	errF4 := exec.Command("xmodmap", "-e", "keycode 70 = NoSymbol").Run()
	if errF4 != nil {
		return errF4
	}

	return nil
}

func (c *UbuntuCommand) EnableAltF4() error {
	errAlt := exec.Command("xmodmap", "-e", "keycode 64 = Alt_L").Run()
	if errAlt != nil {
		return errAlt
	}
	errF4 := exec.Command("xmodmap", "-e", "keycode 70 = F4").Run()
	if errF4 != nil {
		return errF4
	}

	return nil
}

func (c *UbuntuCommand) ChangeUser(user string, password string) error {
	err := exec.Command("sudo", "usermod", "-l", user, password).Run()
	if err != nil {
		return fmt.Errorf("failed to change username: %v", err)
	}

	// Change password
	passwdCmd := exec.Command("sudo", "passwd", user)
	passwdCmd.Stdin = exec.Command("echo", fmt.Sprintf("%s\n%s", password, password)).Stdin
	err = passwdCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to change password: %v", err)
	}

	return nil
}

func (c *UbuntuCommand) SetNumLock(enabled bool) error {
	state := "off"
	if enabled {
		state = "on"
	}

	cmd := exec.Command("numlockx", state)

	// Set DISPLAY to :0 (default X session)
	cmd.Env = append(os.Environ(), "DISPLAY=:0")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set NumLock %s: %v", state, err)
	}

	return nil
}
