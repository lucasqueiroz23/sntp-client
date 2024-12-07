package clientSocket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sntp-client/error-handling"
	"time"
)

const timeout = 20 // if there's been 20 seconds without an answer, we'll try again

type NtpPacket struct {
	Li_vn_mode     uint8
	Stratum        uint8
	Poll           uint8
	Precision      uint8
	RootDelay      uint32
	RootDispersion uint32
	RefId          uint32
	RefTm_s        uint32
	RefTm_f        uint32
	OrigTm_s       uint32
	OrigTm_f       uint32
	RxTm_s         uint32
	RxTm_f         uint32
	TxTm_s         uint32
	TxTm_f         uint32
}

func buildPacketByteArray() []byte {
	packet := new(NtpPacket)
	packet.Li_vn_mode = 0x1B

	message := bytes.NewBuffer(make([]byte, 0, 48))

	writeErr := binary.Write(message, binary.BigEndian, packet)
	if writeErr != nil {
		errorHandling.LogErrorAndExit(writeErr)
	}

	return message.Bytes()

}

func MakeRequest(IPAddress string) []byte {

	// the default port for the NTP protocol
	ntpPort := "123"
	fullAddress := IPAddress + ":" + ntpPort

	conn, dialErr := net.Dial("udp", fullAddress)

	if dialErr != nil {
		errorHandling.LogErrorAndExit(dialErr)
	}

	message := buildPacketByteArray()
	_, writeErr := conn.Write(message)

	if writeErr != nil {
		errorHandling.LogErrorAndExit(writeErr)
	}

	response, readErr := readResponse(conn)

	if readErr != nil {
		errorHandling.LogErrorAndExit(readErr)
	}

	conn.Close()

	return response
}

func readResponse(conn net.Conn) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(timeout * time.Second))

	response := make([]byte, 48) // the size of a packet is 48 bytes
	_, readErr := conn.Read(response)

	if readErr != nil {
		if netErr, ok := readErr.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Timeout de 20 segundos.")
			fmt.Println("Tentando novamente...")
			conn.SetReadDeadline(time.Now().Add(timeout * time.Second))
			_, retryErr := conn.Read(response)
			if retryErr != nil {
				return nil, errors.New("Data/hora: não foi possível contactar servidor")
			}
		}
	}

	return response, nil
}
