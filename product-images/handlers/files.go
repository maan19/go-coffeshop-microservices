package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

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

// UploadREST for RESTful upload
func (f *Files) UploadREST(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "path", r.URL.Path, "id", id, "filename", fn)

	//No need to check id and filename parameters for validation as it
	//is already done by mux router

	f.saveFile(id, fn, rw, r.Body)

}

// Multipart upload
func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Error parsing multipart form", "error", err)
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	f.log.Info("got id:", id)
	if err != nil {
		f.log.Error("Bad id format", "error", err)
		http.Error(rw, "Bad id format", http.StatusBadRequest)
		return
	}

	ff, mh, err := r.FormFile("file")
	if err != nil {
		f.log.Error("Error reading file", "error", err)
		http.Error(rw, "Error reading file", http.StatusBadRequest)
		return
	}

	f.saveFile(r.FormValue("id"), mh.Filename, rw, ff)
}

func (f *Files) saveFile(id string, filename string, rw http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file", "id", id, "filename", filename)

	fp := filepath.Join(id, filename)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Failed to save file", "id", id, "filename", filename, "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
		return
	}
}
