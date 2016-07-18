package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/as27/ranssafe/fileinfo"
	"github.com/gorilla/mux"
)

// GetFileInfo responses with fileinfos to a package
func GetFileInfo(w http.ResponseWriter, r *http.Request) {
	var fileInfos []fileinfo.File
	vars := mux.Vars(r)
	pack := vars["package"]

	rootPath := filepath.Join(
		filepath.ToSlash(Conf.ServerBackupFolder),
		pack)
	log.Println("Walk for: " + rootPath)
	_, err := os.Stat(rootPath)
	if os.IsNotExist(err) == false {
		filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			f, err := fileinfo.New(path)
			if err != nil {
				log.Fatal(err)
			}
			fileInfos = append(fileInfos, f)

			return nil
		})
	}
	b, err := json.Marshal(fileInfos)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}

// PushFile puts the file to the server
func PushFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pack := vars["package"]
	filep := vars["filep"]
	fp := filepath.Join(
		filepath.ToSlash(Conf.ServerBackupFolder),
		pack,
		filep)
	os.MkdirAll(filepath.Dir(fp), 0777)
	f, err := os.Create(fp)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(f, r.Body)
	f.Close()
	r.Body.Close()
	timets, _ := time.Parse(fileinfo.TimestampLayout, r.FormValue("timestamp"))
	err = os.Chtimes(fp, timets, timets)
	if err != nil {
		log.Fatal(err)
	}
}

// GetServerBase returns the basePath from the server that the client is able
// to build a rel path
func GetServerBase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pack := vars["package"]
	fmt.Fprintf(w, "%s/%s", filepath.ToSlash(Conf.ServerBackupFolder), pack)
}
