package main

import (
	"fmt"
	"sntp-client/client-socket"
	"sntp-client/command-line"
	"sntp-client/parser"
)

func main() {
	ipAddress := commandLine.GetIPAddress()
	response := clientSocket.MakeRequest(ipAddress)
	fmt.Println(parser.GetDate(response))
}
