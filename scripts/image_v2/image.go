package image_v2

import (
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nfnt/resize"
)

type Response struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Path    string `json:"path"`
}

func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("func Upload")
	var response Response

	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("image")
	if err != nil {
		response.Status = 3
		response.Message = "Error Uploaded File"
		response.Error = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Detect file type
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		response.Status = 4
		response.Message = "Error Reading File"
		response.Error = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	fileType := http.DetectContentType(fileHeader)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
		response.Status = 5
		response.Message = "Unsupported File Type"
		response.Error = "Only JPEG, PNG and GIF are supported"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Reset file pointer to start
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		response.Status = 6
		response.Message = "Error Resetting File Pointer"
		response.Error = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Decode image
	img, imgType, err := image.Decode(file)
	if err != nil {
		response.Status = 2
		response.Message = "Error Decoding Image"
		response.Error = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Resize image to width 500 preserving the aspect ratio
	resizedImg := resize.Resize(500, 0, img, resize.Lanczos3)

	// Create file
	dstPath := fmt.Sprintf("./images/%s.%s", timestamp, imgType)
	dst, err := os.Create(dstPath)
	if err != nil {
		response.Status = 2
		response.Message = "Error Creating File"
		response.Error = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer dst.Close()

	// Encode and save the resized image
	switch imgType {
	case "jpeg":
		err = jpeg.Encode(dst, resizedImg, nil)
	case "png":
		err = png.Encode(dst, resizedImg)
	case "gif":
		err = gif.Encode(dst, resizedImg, nil)
	}

	if err != nil {
		response.Status = 1
		response.Message = "Error Saving Image"
		response.Error = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = 0
	response.Message = "Successfully Uploaded and Resized File"
	response.Path = os.Getenv("IMAGE_URL") + "/images/" + timestamp + "." + imgType

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
