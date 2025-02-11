package connection

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
)

const HandshakeID = "agentofthings/handshake/0.0.1"

func HandleHandshake(s network.Stream) {
	fmt.Println("Identified a new user, handshaking!")

	defer s.Close()

	buf := make([]byte, 1024)
	n, err := s.Read(buf)
	if err != nil {
		fmt.Println("Error reading from stream:", err)
		return
	}

	fmt.Printf("Received message: %s\n", string(buf[:n]))

	s.Write([]byte("ACK: Connection Accepted"))
}
