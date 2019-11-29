package main

import (
	"fmt"
	"os"

	"github.com/disposedtrolley/crate/internal/client"
	"github.com/disposedtrolley/crate/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no command supplied")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		server.Start()
	case "client":
		client.Start()
	default:
		fmt.Println("`server` or `client` commands are supported")
		os.Exit(1)
	}
}
