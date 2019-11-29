package main

import (
	"fmt"
	"os"

	"github.com/disposedtrolley/crate/internal/client"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("watch directory must be supplied")
		os.Exit(1)
	}
	watchDir := os.Args[1]
	client.Start(watchDir)

}
