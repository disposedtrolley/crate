package client

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"github.com/disposedtrolley/crate/internal/pkg/syncable"
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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				onEvent(watcher, event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(absWatchDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func onEvent(w *fsnotify.Watcher, e fsnotify.Event) {
	// Ignore .swp files
	//if strings.HasSuffix(e.Name, ".swp") {
	//		return
	//}

	// Listen for:
	//  - file created/modified/deleted
	//  - dir created/modified/deleted

	s, err := syncable.NewSyncable(e)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(s)

	if s.FileType() == syncable.FileTypeDir && s.IsCreated() {
		onDirCreate(w, s.Path())
	}

	if s.FileType() == syncable.FileTypeDir && s.IsDeleted() {
		onDirDelete(w, s.Path())
	}
}

func onDirCreate(w *fsnotify.Watcher, dirpath string) error {
	return w.Add(dirpath)
}

func onDirDelete(w *fsnotify.Watcher, dirpath string) error {
	return w.Remove(dirpath)
}
