package main

import (
	"errors"
	"fmt"
	"os"
)

type ntpPacket struct {
	uli_vn_mode     int8
	ustratum        int8
	upoll           int8
	uprecision      int8
	urootDelay      int32
	urootDispersion int32
	urefId          int32
	urefTm_s        int32
	urefTm_f        int32
	uorigTm_s       int32
	uorigTm_f       int32
	urxTm_s         int32
	urxTm_f         int32
	utxTm_s         int32
	txTm_f          uint32
}

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
