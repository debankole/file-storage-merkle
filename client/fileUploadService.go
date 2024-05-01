package main

import (
	"fmt"
	"os"
	"path"

	"github.com/vitaliy/file-storage/common/merkleTree"
)

type FileUploadService struct {
	client *FileServerClient
}

func NewFileUploadService() *FileUploadService {
	return &FileUploadService{client: NewFileServerClient()}
}

func (f *FileUploadService) UploadFiles(dir string) (string, error) {
	key, err := f.client.UploadFiles(dir)
	if err != nil {
		return "", err
	}

	hashes, err := f.GetDirFilesHashes(dir)
	if err != nil {
		return "", err
	}

	merkleRoot, err := merkleTree.GetMerkleRoot(hashes)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(path.Join("merkle_roots", key), os.ModePerm)
	if err != nil {
		return "", err
	}

	file, err := os.Create(path.Join("merkle_roots", key, "merkle_root"))
	defer file.Close()
	if err != nil {
		return "", err
	}

	file.Write(merkleRoot)

	os.RemoveAll(dir)

	return key, nil
}

func (f *FileUploadService) GetFile(key string, num int) ([]byte, string, error) {
	merkleRoot, err := os.ReadFile(path.Join("merkle_roots", key, "merkle_root"))
	if err != nil {
		return nil, "", err
	}

	file, name, err := f.client.GetFile(key, num)
	if err != nil {
		return nil, "", err
	}

	proof, err := f.client.GetProof(key, num)
	if err != nil {
		return nil, "", err
	}

	hash, err := merkleTree.GetHashFromBytes(file)
	if err != nil {
		return nil, "", err
	}
	verificationResult, err := merkleTree.VerifyProof(merkleRoot, num, hash, proof)
	if err != nil {
		return nil, "", err
	}
	if !verificationResult {
		return nil, "", fmt.Errorf("verification failed for file %v", num)
	}

	err = os.WriteFile(fmt.Sprintf("downloads/%v", name), file, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return file, name, nil
}

func (f *FileUploadService) GetDirFilesHashes(dir string) ([][]byte, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	hashes := make([][]byte, 0, len(files))

	for _, file := range files {
		bytes, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		hash, err := merkleTree.GetHashFromBytes(bytes)
		if err != nil {
			return nil, err
		}

		hashes = append(hashes, hash)
	}

	return hashes, nil
}
