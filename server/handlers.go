package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// GetFileInfo responses with fileinfos to a package
func GetFileInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pack := vars["package"]
	fmt.Println(pack)
}
