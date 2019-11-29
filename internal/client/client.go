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
	absWatchDir := validateWatchDir(watchDir)
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

	err = initWatcher(watcher, absWatchDir)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func validateWatchDir(watchDir string) (absWatchDir string) {
	absWatchDir, err := filepath.Abs(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = os.Stat(absWatchDir); os.IsNotExist(err) {
		fmt.Printf("watch directory %s doesn't exist\n", absWatchDir)
		os.Exit(1)
	}

	return absWatchDir
}

func initWatcher(w *fsnotify.Watcher, rootDir string) error {
	err := w.Add(rootDir)
	if err != nil {
		return err
	}

	// walk all nested dirs and add them to the watcher
	err = filepath.Walk(rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsDir() {
				return w.Add(path)
			}
			return nil
		})

	return err
}

// onEvent handles changes to any file or directory being watched by the
// watcher. Roughly, it:
//
// 1) Transforms the filesystem event into a Syncable, which provides
//    some convenience methods to return metadata about the resource.
// 2) Checks if a new directory has been created or a directory has
//    been removed, and calls the event handlers for those.
// 3) Adds the Syncable to a SyncQueue, which is responsible for deciding
//    when to sync the files.
func onEvent(w *fsnotify.Watcher, e fsnotify.Event) {
	s, err := syncable.NewSyncable(e)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(s)

	if s.FileType() == syncable.FileTypeDir && s.IsCreated() {
		err := onDirCreate(w, s.Path())
		if err != nil {
			log.Printf("watch new dir: %v\n", err)
		}
	}

	if s.FileType() == syncable.FileTypeDir && s.IsDeleted() {
		err := onDirDelete(w, s.Path())
		if err != nil {
			log.Printf("unwatch deleted dir: %v\n", err)
		}
	}
}

func onDirCreate(w *fsnotify.Watcher, dirpath string) error {
	return w.Add(dirpath)
}

func onDirDelete(w *fsnotify.Watcher, dirpath string) error {
	return w.Remove(dirpath)
}
