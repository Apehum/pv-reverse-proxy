package repo

import (
	"errors"
	"github.com/BurntSushi/toml"
	"net"
)

type TomlRepository struct {
	servers map[string]*net.UDPAddr
}

type tomlConfig struct {
	Servers map[string]string `toml:"servers"`
}

func NewTomlRepository(filePath string) (*TomlRepository, error) {
	var config tomlConfig
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}

	var servers = make(map[string]*net.UDPAddr)
	for serverHostname, serverIp := range config.Servers {
		serverAddr, err := net.ResolveUDPAddr("udp", serverIp)
		if err != nil {
			return nil, err
		}

		servers[serverHostname] = serverAddr
	}

	return &TomlRepository{servers}, nil
}

func (repo *TomlRepository) GetServerAddress(serverIp string, serverPort uint16) (*net.UDPAddr, error) {
	server := repo.servers[serverIp]
	if server == nil {
		return nil, errors.New("server not found for ip " + serverIp)
	}

	return server, nil
}
