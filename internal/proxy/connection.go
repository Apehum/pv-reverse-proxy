package proxy

import (
	"errors"
	"fmt"
	"log"
	"net"
	packet2 "pv-reverse-proxy/internal/packet"
	"time"
)

type VoiceProxyConnection struct {
	proxy      *VoiceProxy
	clientAddr *net.UDPAddr
	serverConn *net.UDPConn

	secret string
}

func newConnection(proxy *VoiceProxy, packet []byte, clientAddr *net.UDPAddr) (*VoiceProxyConnection, error) {
	voicePacket, err := packet2.DecodePacket(packet)
	if err != nil {
		return nil, err
	}

	if voicePacket.Type != 0x1 {
		return nil, fmt.Errorf("invalid packet type: 0x%x", voicePacket.Type)
	}

	pingPacket, err := packet2.DecodePingPacket(voicePacket)
	if err != nil {
		return nil, err
	}

	if pingPacket.ServerIp == nil {
		return nil, fmt.Errorf("invalid packet server ip; client using aold version of the Plasmo Voice")
	}

	serverAddr, err := proxy.repo.GetServerAddress(*pingPacket.ServerIp, *pingPacket.ServerPort)
	if err != nil {
		return nil, err
	}

	serverConn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, err
	}

	log.Printf("New connection %s; server: %s; server hostname: %s", clientAddr, serverAddr, *pingPacket.ServerIp)

	return &VoiceProxyConnection{
		proxy:      proxy,
		clientAddr: clientAddr,
		serverConn: serverConn,
		secret:     voicePacket.Secret,
	}, nil
}

func (conn *VoiceProxyConnection) writeToServer(packet []byte) error {
	_, err := conn.serverConn.Write(packet)
	return err
}

func (conn *VoiceProxyConnection) writeToClient(packet []byte) error {
	_, err := conn.proxy.conn.WriteToUDP(packet, conn.clientAddr)
	return err
}

func (conn *VoiceProxyConnection) listen() {
	var buffer [1500]byte

	for {
		conn.serverConn.SetReadDeadline(time.Now().Add(30 * time.Second))
		n, err := conn.serverConn.Read(buffer[0:])

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			conn.proxy.mutex.Lock()

			delete(conn.proxy.connections, conn.clientAddr.String())
			conn.serverConn.Close()

			conn.proxy.mutex.Unlock()

			log.Printf("Client %s timed out", conn.clientAddr.String())
			break
		}

		if err != nil {
			log.Println(err)
			continue
		}

		packet := buffer[:n]

		if err = conn.writeToClient(packet); err != nil {
			log.Println(err)
		}
	}
}
