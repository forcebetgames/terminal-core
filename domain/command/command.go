package command

type Command interface {
	ChangeUser(user string, password string) error
	DisableAltF4() error
	EnableAltF4() error
	SetNumLock(enabled bool) error
}

type CommandType string

const (
	COMMAND_DISABLE_ALT_F4 = "DISABLE_ALT_F4"
	COMMAND_ENABLE_ALT_F4  = "ENABLE_ALT_F4"

	COMMAND_DISABLE_KEYS = "DISABLE_KEYS"
	COMMAND_ENABLED_KEYS = "ENABLED_KEYS"
)

func (c CommandType) IsValid() bool {
	switch c {
	case COMMAND_DISABLE_ALT_F4, COMMAND_ENABLE_ALT_F4, COMMAND_DISABLE_KEYS, COMMAND_ENABLED_KEYS:
		return true
	}

	return false
}
