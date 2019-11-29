package client

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
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
	// Listen for:
	//  - file created/modified/deleted
	//  - dir created/modified/deleted

	s, err := NewSyncable(e)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(s)

}

func onDirCreate(w *fsnotify.Watcher, dirpath string) error {
	return w.Add(dirpath)
}

type fileType string

const (
	fileTypeDir     fileType = "DIR"
	fileTypeFile             = "FILE"
	fileTypeDeleted          = "DELETED"
)

type Syncable struct {
	absPath   string
	ftype     fileType
	op        fsnotify.Op
	timestamp time.Time
}

func (s *Syncable) Path() string {
	return s.absPath
}

func (s *Syncable) FileType() fileType {
	return s.ftype
}

func (s *Syncable) Op() fsnotify.Op {
	return s.op
}

func (s *Syncable) Time() time.Time {
	return s.timestamp
}

func (s *Syncable) String() string {
	return fmt.Sprintf("\n====\nresource: %s \ntype: %s \nop: %s \ntime: %s\n====\n", s.Path(), s.FileType(), s.Op(), s.Time().Format("20060102-15:04:05.000"))
}

func (s *Syncable) IsDeleted() bool {
	return s.op&fsnotify.Remove == fsnotify.Remove
}

func NewSyncable(e fsnotify.Event) (*Syncable, error) {
	abspath, err := filepath.Abs(e.Name)
	if err != nil {
		return nil, err
	}

	var ftype fileType
	deleted := e.Op&fsnotify.Remove == fsnotify.Remove
	if deleted {
		ftype = fileTypeDeleted
	} else {
		fi, err := os.Stat(abspath)
		if err != nil {
			return nil, err
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			ftype = fileTypeDir
		case mode.IsRegular():
			ftype = fileTypeFile
		}
	}

	return &Syncable{
		absPath:   abspath,
		ftype:     ftype,
		op:        e.Op,
		timestamp: time.Now(),
	}, nil
}
