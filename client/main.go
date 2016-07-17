package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/as27/ranssafe/distsync"
	"github.com/as27/ranssafe/fileinfo"
)

func main() {
	for _, pack := range Conf.Packages {
		syncer := NewSyncer(Conf.ServerURL, pack.Path)
		syncer.ServerURL = Conf.ServerURL + "/" + pack.Name
		syncer.files = []fileinfo.File{}
		filepath.Walk(pack.Path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			err = syncer.AddFile(path)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		distsync.Distsync(syncer)
	}
}
