package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"path/filepath"

	"strconv"

	"github.com/as27/ranssafe/fileinfo"
)

// ServerURL is the default value for every new syncer
//var ServerURL = "http://localhost:1234"

// ServerFileInfoPath describes the relative URL path to the servers
// API which returns the Fileinfos of a package
var ServerFileInfoPath = "/fileinfo"

// ServerPushPath is the relative path, when pushing a file to the server
var ServerPushPath = "/push"

// ServerBase is the path when the server returns the base file from the conf
var ServerBase = "/serverbase"

// Syncer is a implementation of the Distsyncer interface
type Syncer struct {
	// ServerAdress is the http adess of the server including
	// the port and the package
	// https://localhost:1234/mypackage
	ServerURL      string
	files          []fileinfo.File
	rootPath       string
	serverRootPath string
	newFileinfo    func(string) (fileinfo.File, error)
	osOpen         func(name string) (*os.File, error)
	client         *http.Client
}

// NewSyncer takes a serverAdress and returns a pointer to a
// syncer
func NewSyncer(serverURL string, rootPath string) *Syncer {
	s := Syncer{
		ServerURL: serverURL,
		rootPath:  filepath.ToSlash(rootPath)}
	s.newFileinfo = fileinfo.New
	s.osOpen = os.Open
	s.client = new(http.Client)
	return &s
}

// AddFile adds one file to the Syncer
func (s *Syncer) AddFile(fp string) error {
	fi, err := s.newFileinfo(fp)
	if err != nil {
		return err
	}
	s.files = append(s.files, fi)
	return nil
}

// GetSrcFileInfo implements the distsync interface
func (s *Syncer) GetSrcFileInfo() []fileinfo.File {
	return s.files
}

// GetDistFileInfo implements the distsync interface
// The request uses the ServerFileInfoPath for requesting the fileInfos
// for a package.
func (s *Syncer) GetDistFileInfo() ([]fileinfo.File, error) {
	res, err := http.Get(s.ServerURL + ServerFileInfoPath)
	if err != nil {
		return nil, err
	}
	rbody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	var fi []fileinfo.File
	err = json.Unmarshal(rbody, &fi)
	if err != nil {
		return nil, err
	}
	// Replace the serverBase folder with the rootPath of the client
	res, err = http.Get(s.ServerURL + ServerBase)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	serverBase := string(b)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	for i := range fi {
		fi[i].FilePath = strings.Replace(fi[i].FilePath, serverBase, s.rootPath, 1)

	}
	return fi, nil
}

// PushFile implements the distsync interface
func (s *Syncer) PushFile(fpath string) error {
	// Make URL
	fi, err := fileinfo.New(fpath)
	if err != nil {
		return err
	}
	u, err := url.Parse(s.ServerURL + ServerPushPath + "/" + s.relPath(fpath))
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Set("timestamp", strconv.FormatInt(fi.Timestamp, 10))
	u.RawQuery = v.Encode()
	//log.Println(u.String())
	// Open the local file
	fileReader, err := os.Open(fpath)
	defer fileReader.Close()
	if err != nil {
		return err
	}
	r, err := http.NewRequest("PUT", u.String(), fileReader)
	if err != nil {
		return err
	}
	_, err = s.client.Do(r)
	if err != nil {
		return err
	}

	return nil
}

// GetFile implements the distsync interface
func (s *Syncer) GetFile(string) error {
	return nil
}

func (s *Syncer) relPath(fp string) string {
	rp, err := filepath.Rel(s.rootPath, fp)
	if err != nil {
		log.Fatalln(err)
	}
	return filepath.ToSlash(rp)
}
