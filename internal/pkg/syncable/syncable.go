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
	FileTypeRemoved          = "REMOVED"
)

// Syncable represents a filesystem event which we're interested in syncing.
// Check IsDeleted() before accessing file metadata.
type Syncable struct {
	absPath string
	ftype   FileType
	op      fsnotify.Op
	finfo   os.FileInfo // nil if op == REMOVED || ftype == DELETED
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

func (s *Syncable) ModTime() (time.Time, error) {
	if s.IsDeleted() {
		return time.Now(), fmt.Errorf("deleted resource has no ModTime")
	}
	return s.finfo.ModTime(), nil
}

func (s *Syncable) String() string {
	return fmt.Sprintf("\n====\nresource: %s \ntype: %s \nop: %s \n====\n", s.Path(), s.FileType(), s.Op())
}

func (s *Syncable) IsDeleted() bool {
	return s.ftype == FileTypeRemoved
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
	var finfo os.FileInfo
	deletedOrRenamed := e.Op&fsnotify.Remove == fsnotify.Remove || e.Op&fsnotify.Rename == fsnotify.Rename
	if deletedOrRenamed {
		ftype = FileTypeRemoved
	} else {
		finfo, err := os.Stat(abspath)
		if err != nil {
			return nil, err
		}

		switch mode := finfo.Mode(); {
		case mode.IsDir():
			ftype = FileTypeDir
		case mode.IsRegular():
			ftype = FileTypeFile
		}
	}

	return &Syncable{
		absPath: abspath,
		ftype:   ftype,
		op:      e.Op,
		finfo:   finfo,
	}, nil
}
