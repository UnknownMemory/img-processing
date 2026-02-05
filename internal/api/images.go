package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/unknownmemory/img-processing/internal/database"
	"github.com/unknownmemory/img-processing/internal/rabbitmq"
	"github.com/unknownmemory/img-processing/internal/shared"
)

func (app *Application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "An error occurred during the upload")
		return
	}

	file, handler, err := r.FormFile("image-upload")
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "An error occurred during the upload")
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
		app.errorResponse(w, r, http.StatusBadRequest, "An error occurred during the upload")
		return
	}
	mimeType := http.DetectContentType(buffer)
	if !strings.HasPrefix(mimeType, "image/") {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid MIME type")
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
	imgUUID, err := q.CreateImage(context.Background(), *data)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "An error occurred during the upload")
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "An error occurred during the upload")
		return
	}

	key := fmt.Sprintf("%d/%s/original", userId, imgUUID)
	url := fmt.Sprintf("%s/%s", app.s3.BucketPublicURL, key)
	_, err = app.s3.Upload(key, file, mimeType)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "An error occurred during the upload")
		return
	}

	fileData := map[string]string{
		"imageID":  imgUUID.String(),
		"filename": handler.Filename,
		"size":     strconv.FormatInt(handler.Size, 10),
		"url":      url,
	}

	err = app.writeJSON(w, http.StatusCreated, fileData, nil)
	if err != nil {
		app.internalErrorResponse(w, r, err)
		return
	}
}

func (app *Application) transform(w http.ResponseWriter, r *http.Request) {
	data := &shared.ImageTransform{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		app.internalErrorResponse(w, r, err)
		return
	}

	imageUUID, err := uuid.Parse(data.ImageID)
	if err != nil {
		app.logger.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid Image ID")
		return
	}

	userId := r.Context().Value("user_id").(int64)
	q := db.New(app.db)
	imageQuery := &db.ImageExistsParams{
		UserID: pgtype.Int8{Int64: userId, Valid: true},
		Uid:    pgtype.UUID{Bytes: imageUUID, Valid: true},
	}

	image, err := q.ImageExists(context.Background(), *imageQuery)
	if err != nil {
		app.internalErrorResponse(w, r, err)
		return
	}

	if !image {
		errImg := map[string]interface{}{"error": "Image not found"}
		err = app.writeJSON(w, http.StatusNotFound, errImg, nil)
		if err != nil {
			app.internalErrorResponse(w, r, err)
			return
		}
	}

	transformParams := &db.CreateTransformParams{
		OriginalImage: pgtype.UUID{Bytes: imageUUID, Valid: true},
		UserID:        pgtype.Int8{Int64: userId, Valid: true},
	}
	transform, err := q.CreateTransform(context.Background(), *transformParams)
	if err != nil {
		app.internalErrorResponse(w, r, err)
		return
	}
	rmq := rabbitmq.NewWorker(os.Getenv("RABBIT_MQ"), app.logger)
	rmq.Send("image", data, strconv.FormatInt(userId, 10))

	err = app.writeJSON(w, http.StatusAccepted, transform, nil)
	if err != nil {
		app.internalErrorResponse(w, r, err)
		return
	}
}

func (app *Application) getImage(w http.ResponseWriter, r *http.Request) {
	uuidParams := r.PathValue("id")
	userId := r.Context().Value("user_id").(int64)

	imgUUID, err := uuid.Parse(uuidParams)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid Image ID")
		return
	}

	q := db.New(app.db)
	imgParams := &db.GetImageParams{
		UserID: pgtype.Int8{Int64: userId, Valid: true},
		Uid:    pgtype.UUID{Bytes: imgUUID, Valid: true},
	}

	image, err := q.GetImage(context.Background(), *imgParams)
	if err != nil {
		app.internalErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, image, nil)
	if err != nil {
		app.internalErrorResponse(w, r, err)
		return
	}
}
