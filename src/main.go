package main

import (
	"errors"
	"fmt"
	"os"
)

func getCommandLineArguments() ([]string, error) {
	args := os.Args[1:]
	correctUsage := "Modo correto de usar:\n./client <ip do servidor NTP>"

	if len(args) != 1 {
		return nil, errors.New(correctUsage)
	}

	return args, nil
}

func main() {

	args, err := getCommandLineArguments()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("args: ")
	for _, arg := range args {
		fmt.Println(arg)
	}
}
