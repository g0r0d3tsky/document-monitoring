package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main-service/internal/usecase/model"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type ContentHandler struct {
	textService    TextContentService
	contentService ContentService
}

type TextContentService interface {
	DeleteContent(ctx context.Context, filename string) error
	GetTextContent(ctx context.Context, filename string) (*model.Content, error)
}

type ContentService interface {
	SendFileToKafka(filename string, fileContent []byte) error
	DeleteContentFromServer(ctx context.Context, filename string) error
	GetFile(ctx context.Context, filename string) ([]byte, error)
}

func NewContentHandler(textService TextContentService, contentService ContentService) *ContentHandler {
	return &ContentHandler{
		textService:    textService,
		contentService: contentService,
	}
}

func (c *ContentHandler) PostContent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "multipart/form-data" {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusUnsupportedMediaType)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // ограничение на 10 MB
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	filename := r.FormValue("filename")

	err = c.contentService.SendFileToKafka(filename, fileContent)
	if err != nil {
		http.Error(w, "Failed to send file to Kafka", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
}

func (c *ContentHandler) DeleteContentName(w http.ResponseWriter, r *http.Request, name string) {
	extension := filepath.Ext(name)

	if extension == ".txt" {
		err := c.textService.DeleteContent(r.Context(), name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete content from MongoDB: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		err := c.contentService.DeleteContentFromServer(r.Context(), name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete content from server: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Content deleted successfully"})
}

// GetContentName - Handles getting content by filename
func (c *ContentHandler) GetContentName(w http.ResponseWriter, r *http.Request, name string) {
	fileExtension := strings.ToLower(strings.TrimPrefix(name, "."))

	var content *model.Content
	var err error

	if fileExtension == "txt" {
		content, err = c.textService.GetTextContent(r.Context(), name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(content.Text))

	} else {
		contentBytes, err := c.contentService.GetFile(r.Context(), name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		contentType := detectContentType(fileExtension)
		w.Header().Set("Content-Type", contentType)

		w.Write(contentBytes)
	}
}

func detectContentType(extension string) string {
	if !filepath.IsAbs(extension) {
		extension = "." + extension
	}

	contentType := mime.TypeByExtension(extension)
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}

func (c *ContentHandler) PutContentName(w http.ResponseWriter, r *http.Request, name string) {
	//TODO: fix that
}
