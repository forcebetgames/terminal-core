package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	debugEnabled bool
	infoLogger   *log.Logger
	errorLogger  *log.Logger
)

// Init initializes the logger with debug mode
func Init(debug bool) {
	debugEnabled = debug

	if debugEnabled {
		// Se debug está habilitado, usa stdout/stderr normalmente
		infoLogger = log.New(os.Stdout, "", 0)
		errorLogger = log.New(os.Stderr, "", 0)
	} else {
		// Se debug está desabilitado, descarta toda saída
		infoLogger = log.New(io.Discard, "", 0)
		errorLogger = log.New(io.Discard, "", 0)
	}
}

// Print imprime mensagem (sem quebra de linha) se debug estiver habilitado
func Print(v ...interface{}) {
	if debugEnabled {
		fmt.Print(v...)
	}
}

// Println imprime mensagem com quebra de linha se debug estiver habilitado
func Println(v ...interface{}) {
	if debugEnabled {
		fmt.Println(v...)
	}
}

// Printf imprime mensagem formatada se debug estiver habilitado
func Printf(format string, v ...interface{}) {
	if debugEnabled {
		fmt.Printf(format, v...)
	}
}

// IsDebug retorna se o modo debug está habilitado
func IsDebug() bool {
	return debugEnabled
}
