package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// FileInfoPath describes the relative URL when the fileinfo slice
// should be returned
var FileInfoPath = "/fileinfo"

// PushPath describes the path when a file is pushed to the server
var PushPath = "/push"

// ServerBase is the path when the server returns the base file from the conf
var ServerBase = "/serverbase"

func main() {
	LoadConf()
	router := mux.NewRouter()
	router.HandleFunc(`/{package}`+ServerBase, GetServerBase).Methods("GET")
	router.HandleFunc(`/{package}`+FileInfoPath, GetFileInfo).Methods("GET")
	router.HandleFunc(`/{package}`+PushPath+`/{filep:.*}`, PushFile).Methods("PUT")
	log.Println("Starting server at port " + Conf.ServerPort)
	log.Println("Server folder: " + Conf.ServerBackupFolder)
	log.Fatal(http.ListenAndServe(Conf.ServerPort, router))
}
