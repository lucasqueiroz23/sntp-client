package main

import (
	"fmt"
	"sntp-client/client-socket"
	"sntp-client/command-line"
)

func main() {

	ipAddress := commandLine.GetIPAddress()
	packet := clientSocket.BuildPacket()

	fmt.Println("ip address: ", ipAddress)
	fmt.Println("packet:", packet)
}
