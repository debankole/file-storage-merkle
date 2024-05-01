package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	fileservice "github.com/vitaliy/file-storage/server/fileService"
	filestore "github.com/vitaliy/file-storage/server/fileStore"
)

func uploadFilesHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileHeaders := r.MultipartForm.File["files"]

	files := make([]filestore.FileInfo, 0, len(fileHeaders))

	for _, file := range fileHeaders {
		f, err := file.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		files = append(files, filestore.FileInfo{R: f, Name: file.Filename})
	}

	key, err := fileservice.NewFileService().StoreFiles(nil, files)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error uploading files", http.StatusInternalServerError)
	}

	response := UploadResponse{
		Key: key,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("filenumber")
	key := r.URL.Query().Get("key")

	numberInt, err := strconv.Atoi(number)
	if err != nil {
		http.Error(w, "Invalid file number", http.StatusBadRequest)
		return
	}

	file, name, err := fileservice.NewFileService().GetFile(key, numberInt)

	if err != nil {
		http.Error(w, "Error getting file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", name))
	io.Copy(w, bytes.NewReader(file))
}

func getProofHandler(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("filenumber")
	key := r.URL.Query().Get("key")

	numberInt, err := strconv.Atoi(number)
	if err != nil {
		http.Error(w, "Invalid file number", http.StatusBadRequest)
		return
	}

	proof, err := fileservice.NewFileService().GetProof(key, numberInt)

	if err != nil {
		http.Error(w, "Error getting proof", http.StatusInternalServerError)
		return
	}
	proofResponse := ProofResponse{}

	for _, v := range proof {
		proofResponse.Proof = append(proofResponse.Proof, hex.EncodeToString(v))
	}
	jsonResponse, err := json.Marshal(proofResponse)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

type UploadResponse struct {
	Key string `json:"key"`
}

type ProofResponse struct {
	Proof []string
}

func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		uploadFilesHandler(w, r)
	})
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		getFileHandler(w, r)
	})
	http.HandleFunc("/proof", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		getProofHandler(w, r)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
