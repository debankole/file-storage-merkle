package fileservice

import (
	"os"
	"strings"
	"testing"

	"github.com/vitaliy/file-storage/common/merkleTree"
	filestore "github.com/vitaliy/file-storage/server/fileStore"
)

func TestStoreFiles(t *testing.T) {

	os.RemoveAll("files")
	file1 := NewFileInfo("test1")
	file2 := NewFileInfo("test2")
	file3 := NewFileInfo("test3")
	file4 := NewFileInfo("test4")
	file5 := NewFileInfo("test5")
	file6 := NewFileInfo("test6")
	file7 := NewFileInfo("test7")

	service := NewFileService()
	key, err := service.StoreFiles(nil, []filestore.FileInfo{*file1, *file2, *file3, *file4, *file5, *file6, *file7})
	if err != nil {
		t.Fatalf("Error storing files: %v", err)
	}

	merkleTreeBytes, err := filestore.NewFileStore().GetFileByName(key, filestore.MerkleTreeFileName)
	if err != nil {
		t.Fatalf("Error getting merkle tree: %v", err)
	}

	var tree merkleTree.MerkleTree
	err = merkleTree.UnmarshalTree(merkleTreeBytes, &tree)
	if err != nil {
		t.Fatalf("Error unmarshalling tree: %v", err)
	}

	verifyFile(service, key, t, tree.Root.Hash, 0, "test1")
	verifyFile(service, key, t, tree.Root.Hash, 1, "test2")
	verifyFile(service, key, t, tree.Root.Hash, 2, "test3")
	verifyFile(service, key, t, tree.Root.Hash, 3, "test4")
	verifyFile(service, key, t, tree.Root.Hash, 4, "test5")
	verifyFile(service, key, t, tree.Root.Hash, 5, "test6")
	verifyFile(service, key, t, tree.Root.Hash, 6, "test7")
}

func NewFileInfo(name string) *filestore.FileInfo {
	return &filestore.FileInfo{Name: name, R: strings.NewReader(name)}
}

func verifyFile(service *FileService, key string, t *testing.T, rootHash []byte, index int, name string) {
	file, actualName, err := service.GetFile(key, index)
	if err != nil {
		t.Fatalf("Error getting file: %v", err)
	}

	proof, err := service.GetProof(key, index)

	if name != actualName {
		t.Fatalf("Name mismatch for file %d", index)
	}

	if err != nil {
		t.Fatalf("Error getting file: %v", err)
	}

	hash, err := merkleTree.GetHashFromBytes(file)
	if err != nil {
		t.Fatalf("Error getting hash: %v", err)
	}

	verificationResult, err := merkleTree.VerifyProof(rootHash, index, hash, proof)
	if err != nil {
		t.Fatalf("Error verifying proof: %v", err)
	}

	if !verificationResult {
		t.Fatalf("Verification failed for file %d", index)
	}
}
