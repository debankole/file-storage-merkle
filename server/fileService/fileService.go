package fileservice

import (
	"github.com/google/uuid"
	merkleTree "github.com/vitaliy/file-storage/common/merkleTree"
	filestore "github.com/vitaliy/file-storage/server/fileStore"
)

type FileService struct {
	store *filestore.FileStore
}

func NewFileService() *FileService {
	return &FileService{store: filestore.NewFileStore()}
}

func (f FileService) StoreFiles(key *string, files []filestore.FileInfo) (string, error) {
	if key == nil {
		newUuid := uuid.New().String()
		key = &newUuid
	}

	hashes, err := f.store.StoreFiles(*key, files)
	if err != nil {
		return "", err
	}

	tree, err := merkleTree.NewMerkleTree(hashes)
	if err != nil {
		return "", err
	}

	treeBytes, err := merkleTree.MarshalTree(tree)
	if err != nil {
		return "", err
	}

	f.store.StoreFile(*key, filestore.MerkleTreeFileName, treeBytes)

	return *key, nil
}

func (f FileService) GetProof(key string, number int) ([][]byte, error) {
	treeBytes, err := f.store.GetFileByName(key, filestore.MerkleTreeFileName)
	if err != nil {
		return nil, err
	}

	tree := &merkleTree.MerkleTree{}
	err = merkleTree.UnmarshalTree(treeBytes, tree)
	if err != nil {
		return nil, err
	}

	proof := merkleTree.GetProof(tree, number)

	return proof, nil
}

func (f FileService) GetFile(key string, number int) ([]byte, string, error) {
	file, name, err := f.store.GetFileByNumber(key, number)
	if err != nil {
		return nil, "", err
	}

	return file, name, nil
}
