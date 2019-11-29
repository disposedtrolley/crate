package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/disposedtrolley/crate/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("port must be supplied")
		os.Exit(1)
	}
	port, _ := strconv.Atoi(os.Args[1])
	server.Start(port)
}
