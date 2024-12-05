package commandLine

import (
	"errors"
	"os"
	"sntp-client/error-handling"
)

func getCommandLineArguments() ([]string, error) {
	args := os.Args[1:]
	correctUsage := "Erro no uso. Modo correto de usar:\n./client <ip do servidor NTP>"

	if len(args) != 1 {
		return nil, errors.New(correctUsage)
	}

	return args, nil
}

func GetIPAddress() string {
	args, err := getCommandLineArguments()

	if err != nil {
		errorHandling.LogErrorAndExit(err)
	}

	return args[0]
}
