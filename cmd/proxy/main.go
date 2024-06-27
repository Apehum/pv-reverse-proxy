package main

import (
	"flag"
	"os"
	"pv-reverse-proxy/internal/proxy"
	"pv-reverse-proxy/internal/repo"
)

func main() {
	var flagHelp = flag.Bool("h", false, "Show help")
	var listenPort = flag.Int("p", 30000, "Proxy listen port")

	flag.Parse()
	if *flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	tomlRepo, err := repo.NewTomlRepository("servers.toml")
	if err != nil {
		panic(err)
	}

	voiceProxy, err := proxy.NewProxy(*listenPort, tomlRepo)
	if err != nil {
		panic(err)
	}

	voiceProxy.Listen()
}
