package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type FileServerClient struct {
}

type UploadResponse struct {
	Key string `json:"key"`
}

type ProofResponse struct {
	Proof []string
}

func NewFileServerClient() *FileServerClient {
	return &FileServerClient{}
}

func (f *FileServerClient) GetFile(key string, num int) ([]byte, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/files?key=%v&filenumber=%v", FileServerUrl, key, num), nil)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		return nil, "", err
	}
	filename := params["filename"]

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return response, filename, nil
}

func (f *FileServerClient) GetProof(key string, num int) ([][]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/proof?key=%v&filenumber=%v", FileServerUrl, key, num), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var proofResponse ProofResponse
	err = json.NewDecoder(resp.Body).Decode(&proofResponse)
	if err != nil {
		return nil, err
	}

	proof := make([][]byte, 0, len(proofResponse.Proof))

	for _, v := range proofResponse.Proof {
		decoded, err := hex.DecodeString(v)
		if err != nil {
			return nil, err
		}
		proof = append(proof, decoded)
	}

	return proof, nil
}

func (f *FileServerClient) UploadFiles(dirName string) (string, error) {
	entries, err := os.ReadDir(dirName)
	if err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	slices.SortFunc(entries, func(a fs.DirEntry, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	for _, entry := range entries {
		err := addFileMultipart(writer, filepath.Join(dirName, entry.Name()))
		if err != nil {
			return "", err
		}
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/upload", FileServerUrl), body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var uploadResponse UploadResponse
	err = json.NewDecoder(resp.Body).Decode(&uploadResponse)
	if err != nil {
		return "", err
	}

	return uploadResponse.Key, nil
}

func addFileMultipart(writer *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	part, err := writer.CreateFormFile("files", filepath.Base(file.Name()))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}
