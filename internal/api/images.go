package api

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/unknownmemory/img-processing/internal/database"
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
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			app.logger.Println(err)
		}
	}(file)

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "An error occurred during the upload", http.StatusBadRequest)
		return
	}
	mimeType := http.DetectContentType(buffer)
	if !strings.HasPrefix(mimeType, "image/") {
		app.logger.Println(err)
		http.Error(w, "Invalid MIME type", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value("user_id").(int64)
	data := &db.CreateImageParams{
		UserID:   pgtype.Int8{Int64: userId, Valid: true},
		Filename: handler.Filename,
		FileSize: pgtype.Int8{Int64: handler.Size, Valid: true},
		Mime:     mimeType,
	}

	q := db.New(app.db)
	uuid, err := q.CreateImage(context.Background(), *data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "An error occurred during the upload", http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "An error occurred during the upload", http.StatusBadRequest)
	}

	key := fmt.Sprintf("%d/%s/original", userId, uuid)
	url := fmt.Sprintf("%s/%s", app.s3.BucketPublicURL, key)
	_, err = app.s3.Upload(key, file, mimeType)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "An error occurred during the upload", http.StatusBadRequest)
		return
	}

	fileData := map[string]string{
		"filename": handler.Filename,
		"size":     strconv.FormatInt(handler.Size, 10),
		"url":      url,
	}

	err = app.writeJSON(w, http.StatusCreated, fileData, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
