package filestore

import (
	"bytes"
	"io"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/vitaliy/file-storage/common/merkleTree"
)

const MerkleTreeFileName = "_merkleTree.json"

const Dir = "files"

func NewFileStore() *FileStore {
	return &FileStore{}
}

type FileStore struct {
}

type FileInfo struct {
	Name string
	R    io.Reader
}

func (f FileStore) StoreFile(key string, name string, content []byte) error {
	return os.WriteFile(path.Join(Dir, key, name), content, os.ModePerm)
}

func (f FileStore) StoreFiles(key string, files []FileInfo) ([][]byte, error) {
	err := cleanupDir(key)
	if err != nil {
		return nil, err
	}

	hashes := make([][]byte, 0, len(files))

	slices.SortFunc(files, func(a FileInfo, b FileInfo) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, f := range files {
		filePath := path.Join(Dir, key, f.Name)
		newFile, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer newFile.Close()

		var buf bytes.Buffer
		tee := io.TeeReader(f.R, &buf)

		hash, err := merkleTree.GetHashFromReader(tee)
		if err != nil {
			return nil, err
		}

		hashes = append(hashes, hash)

		bytes.NewReader(buf.Bytes())
		_, err = io.Copy(newFile, bytes.NewReader(buf.Bytes()))
		if err != nil {
			return nil, err
		}
	}

	return hashes, nil
}

func (f FileStore) GetFileByNumber(key string, number int) ([]byte, string, error) {
	fileNames, err := os.ReadDir(path.Join(Dir, key))
	if err != nil {
		return nil, "", err
	}

	for i, file := range fileNames {
		if file.Name() == MerkleTreeFileName {
			fileNames = append(fileNames[:i], fileNames[i+1:]...)
			break
		}
	}

	slices.SortFunc(fileNames, func(a os.DirEntry, b os.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	file, err := os.ReadFile(path.Join(Dir, key, fileNames[number].Name()))

	return file, fileNames[number].Name(), err
}

func (f FileStore) GetFileByName(key string, name string) ([]byte, error) {
	return os.ReadFile(path.Join(Dir, key, name))
}

func cleanupDir(key string) error {
	err := os.RemoveAll(path.Join(Dir, key))
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(Dir, key), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
