package api

import (
	"net/http"
	"strconv"
	"strings"
)

func (app *Application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "An error occurred during the upload", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image-upload")
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "An error occurred during the upload", http.StatusBadRequest)
		return
	}
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "An error occurred during the upload", http.StatusBadRequest)
		return
	}
	mimeType := http.DetectContentType(buffer)
	if !strings.Contains(mimeType, "image/") {
		app.logger.Println(err)
		http.Error(w, "Invalid MIME type", http.StatusBadRequest)
		return
	}

	fileData := map[string]string{"filename": handler.Filename, "size": strconv.FormatInt(handler.Size, 10)}

	err = app.writeJSON(w, http.StatusOK, fileData, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
