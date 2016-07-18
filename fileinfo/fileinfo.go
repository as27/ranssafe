// Package fileinfo is for the exange of metadata of each file. The package is
// used by server and client.
package fileinfo

import (
	"os"
	"path/filepath"
	"strconv"
)

// File is the metadata of a file, which is used to compare files before copying.
type File struct {
	//FilePath is the path to the file starting at the root folder of each package
	FilePath string `json:"filePath"`
	//Timestamp represents the last modification time
	Timestamp int64 `json:"timestamp"`
}

// TimestampLayout defines the timestamps
const TimestampLayout string = "20060102150405"

//New returns a new instance of a File struct.
func New(fp string) (File, error) {
	var f File
	f.FilePath = filepath.ToSlash(fp)
	var err error
	f.Timestamp, err = getTimestamp(fp)
	return f, err
}

func getTimestamp(fp string) (int64, error) {
	fi, err := os.Stat(fp)
	if err != nil {
		return 0, err
	}
	lastMod := fi.ModTime()
	return strconv.ParseInt(lastMod.Format(TimestampLayout), 10, 64)
}

// The Encoder interface represents a Encode method of any encoding library
// the exchange format is open at this point
type Encoder interface {
	Encode(interface{}) error
}

// The Decoder interface represents a Decode method analog to the Encoder
type Decoder interface {
	Decode(interface{}) error
}
