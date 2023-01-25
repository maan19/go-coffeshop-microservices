package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/maan19/go-coffeshop-microservices/product-images/files"
)

// Files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new Files handler
func NewFiles(l hclog.Logger, s files.Storage) *Files {
	return &Files{
		log:   l,
		store: s,
	}
}

// ServeHTTP implements the http.Handler interface
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "path", r.URL.Path, "id", id, "filename", fn)

	//No need to check id and filename parameters for validation as it
	//is already done by mux router

	f.saveFile(id, fn, rw, r)

}

func (f *Files) saveFile(id string, filename string, rw http.ResponseWriter, r *http.Request) {
	f.log.Info("Save file", "id", id, "filename", filename)

	fp := filepath.Join(id, filename)
	err := f.store.Save(fp, r.Body)
	if err != nil {
		f.log.Error("Failed to save file", "id", id, "filename", filename, "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
		return
	}
}
