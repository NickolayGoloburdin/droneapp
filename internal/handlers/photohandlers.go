package handlers

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	repository "github.com/Nickolaygoloburdin/droneapp/internal/database"
	"github.com/julienschmidt/httprouter"
)

type PhotoHandler struct {
	ctx  context.Context
	repo *repository.Repository
}

func NewPhotoHandler(ctx context.Context, repo *repository.Repository) *PhotoHandler {
	return &PhotoHandler{ctx, repo}
}

func (photo PhotoHandler) PostFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	path := FileSave(r)
	if path == "" {
		http.Error(w, "emtypath", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func FileSave(r *http.Request) string {
	err := r.ParseMultipartForm(8024000)
	if err != nil {
		return ""
	}
	n := r.Form.Get("name")
	// Retrieve the file from form data
	f, h, err := r.FormFile("file")
	if err != nil {
		return ""
	}
	defer f.Close()
	path := filepath.Join("/home/nickolay/droneapp/", "files")
	_ = os.MkdirAll(path, os.ModePerm)
	fullPath := path + "/" + n + ".jpg"
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return ""
	}
	defer file.Close()
	// Copy the file to the destination path
	_, err = io.Copy(file, f)
	if err != nil {
		return ""
	}
	return n + filepath.Ext(h.Filename)

}
