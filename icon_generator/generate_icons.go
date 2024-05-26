package icongenerator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"forger/model"
	"image"
	"net/http"
)

func BuildIcon(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data, including the image
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	// Get the file from the form data
	file, _, err := r.FormFile("image") // Ignore the second return value (handler)
	if err != nil {
		http.Error(w, "Error retrieving file from form data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	srcImage, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Error decoding image", http.StatusInternalServerError)
		return
	}

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	IOSmageResizer(zipWriter, srcImage, model.IOSResizeMetaList, IOS)
	IOSmageResizer(zipWriter, srcImage, model.AndroidResizeMetaList, Android)

	err = zipWriter.Close()
	if err != nil {
		http.Error(w, "Error closing zip writer", http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"images.zip\"")
	w.Header().Set("Content-Length", fmt.Sprint(zipBuffer.Len()))

	// Write the zip buffer to the response
	_, err = w.Write(zipBuffer.Bytes())
	if err != nil {
		http.Error(w, "Error writing zip content to response", http.StatusInternalServerError)
		return
	}

}
