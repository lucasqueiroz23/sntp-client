package main

import (
	"fmt"
	"os"
	"sntp-client/cli"
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

func buildPacket() *ntpPacket {
	packet := new(ntpPacket)
	packet.uli_vn_mode = 0x1B
	return packet
}

func checkIPAddress() bool {
	return true
}

func main() {

	args, err := cli.GetCommandLineArguments()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	packet := buildPacket()

	fmt.Println("args: ")
	for _, arg := range args {
		fmt.Println(arg)
	}

	fmt.Println("packet:")
	fmt.Println(packet)
}
