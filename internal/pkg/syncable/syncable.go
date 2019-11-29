package syncable

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileType string

const (
	FileTypeDir     FileType = "DIR"
	FileTypeFile             = "FILE"
	FileTypeDeleted          = "DELETED"
)

type Syncable struct {
	absPath   string
	ftype     FileType
	op        fsnotify.Op
	timestamp time.Time
}

func (s *Syncable) Path() string {
	return s.absPath
}

func (s *Syncable) FileType() FileType {
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
	return s.ftype == FileTypeDeleted
}

func (s *Syncable) IsCreated() bool {
	return s.op&fsnotify.Create == fsnotify.Create
}

func NewSyncable(e fsnotify.Event) (*Syncable, error) {
	abspath, err := filepath.Abs(e.Name)
	if err != nil {
		return nil, err
	}

	var ftype FileType
	deleted := e.Op&fsnotify.Remove == fsnotify.Remove
	if deleted {
		ftype = FileTypeDeleted
	} else {
		fi, err := os.Stat(abspath)
		if err != nil {
			return nil, err
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			ftype = FileTypeDir
		case mode.IsRegular():
			ftype = FileTypeFile
		}
	}

	return &Syncable{
		absPath:   abspath,
		ftype:     ftype,
		op:        e.Op,
		timestamp: time.Now(),
	}, nil
}
