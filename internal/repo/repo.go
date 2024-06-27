package repo

import "net"

type ServerRepository interface {
	GetServerAddress(serverIp string, serverPort uint16) (*net.UDPAddr, error)
}
