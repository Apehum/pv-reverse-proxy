package proxy

import (
	"fmt"
	"log"
	"net"
	"pv-reverse-proxy/internal/repo"
	"sync"
)

type VoiceProxy struct {
	conn        *net.UDPConn
	connections map[string]*VoiceProxyConnection
	repo        repo.ServerRepository

	mutex sync.Mutex
}

func NewProxy(port int, repo repo.ServerRepository) (*VoiceProxy, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	log.Println("Listening on", addr)

	return &VoiceProxy{
		conn:        conn,
		connections: map[string]*VoiceProxyConnection{},
		repo:        repo,
	}, nil
}

func (proxy *VoiceProxy) Listen() {
	var buffer [1500]byte

	for {
		n, clientAddr, err := proxy.conn.ReadFromUDP(buffer[:])
		if err != nil {
			log.Println(err)
			continue
		}

		packet := buffer[:n]
		clientAddrStr := clientAddr.String()

		proxy.mutex.Lock()
		// todo: maybe it's better to use client's secret for identification, instead of address
		conn, ok := proxy.connections[clientAddrStr]
		if !ok {
			conn, err = newConnection(proxy, packet, clientAddr)
			if err != nil {
				log.Println(err)
				proxy.mutex.Unlock()
				continue
			}

			proxy.connections[clientAddrStr] = conn

			go conn.listen()
		}

		proxy.mutex.Unlock()
		if err = conn.writeToServer(packet); err != nil {
			log.Println(err)
		}
	}
}
