package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/as27/ranssafe/distsync"
	"github.com/as27/ranssafe/fileinfo"
)

func main() {
	for _, pack := range Conf.Packages {
		syncer := makeSyncer(pack)
		distsync.Distsync(syncer)
	}
}

func makeSyncer(pack Pack) *Syncer {
	syncer := NewSyncer(Conf.ServerURL, pack.Path)
	syncer.ServerURL = Conf.ServerURL + "/" + pack.Name
	syncer.files = []fileinfo.File{}
	filepath.Walk(pack.Path, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			if skipDir(path) {
				log.Println("SkipDir: " + path)
				return filepath.SkipDir
			}
			return nil
		}
		err = syncer.AddFile(path)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	return syncer
}

func skipDir(p string) bool {
	for _, d := range Conf.SkipDir {
		if strings.Contains(p, d) {
			return true
		}
	}
	return false
}
