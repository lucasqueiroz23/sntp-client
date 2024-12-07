// Package commandLine checks if the IP Address or
// hostname of the NTP server is given in the command line.
package commandLine

import (
	"errors"
	"os"
	"sntp-client/error-handling"
)

// Checks if only the address of the NTP server is given
// and returns it
func checkCommandLineInput() ([]string, error) {
	args := os.Args[1:]
	correctUsage := "Erro no uso. Modo correto de usar:\n./client <ip do servidor NTP>"

	if len(args) != 1 {
		return nil, errors.New(correctUsage)
	}

	return args, nil
}

// Returns the IP Address or hostname if it's given as an argument
func GetIPAddress() string {
	args, err := checkCommandLineInput()

	if err != nil {
		errorHandling.LogErrorAndExit(err)
	}

	return args[0]
}
