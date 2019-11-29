package client

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Start(watchDir string) {
	absWatchDir, err := filepath.Abs(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = os.Stat(absWatchDir); os.IsNotExist(err) {
		fmt.Printf("watch directory %s doesn't exist\n", absWatchDir)
		os.Exit(1)
	}

	fmt.Printf("client started watching %s\n", absWatchDir)
}
