package main

import (
	"sntp-client/client-socket"
	"sntp-client/command-line"
)

func main() {

	ipAddress := commandLine.GetIPAddress()
	clientSocket.MakeRequest(ipAddress)

}
