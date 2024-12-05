package main

import (
	"fmt"
	"sntp-client/cli"
	"sntp-client/client-socket"
)

func main() {

	ipAddress := cli.GetIPAddress()
	packet := clientSocket.BuildPacket()

	fmt.Println("ip address: ", ipAddress)
	fmt.Println("packet:", packet)
}
