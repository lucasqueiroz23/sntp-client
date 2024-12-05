package cli

import (
	"errors"
	"os"
)

func GetCommandLineArguments() ([]string, error) {
	args := os.Args[1:]
	correctUsage := "Erro no uso. Modo correto de usar:\n./client <ip do servidor NTP>"

	if len(args) != 1 {
		return nil, errors.New(correctUsage)
	}

	return args, nil
}
