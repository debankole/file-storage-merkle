package merkleTree

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io"
	"slices"
)

func GetProof(tree *MerkleTree, index int) [][]byte {

	proof := make([][]byte, 0)

	node := tree.Root
	numLeafs := *node.Num/2 + 1

	for node.Left != nil {
		if index >= numLeafs/2 {
			proof = append(proof, node.Left.Hash)
			node = node.Right
			index = index - numLeafs/2
		} else {
			proof = append(proof, node.Right.Hash)
			node = node.Left
		}
		numLeafs = numLeafs / 2
	}

	slices.Reverse(proof)

	return proof
}

// func GetProofNodes(tree *MerkleTree, index int) []Node {

// 	proof := make([][]byte, 0)
// 	nodes := make([]Node, 0)

// 	node := tree.Root
// 	numLeafs := *node.Num/2 + 1

// 	for node.Left != nil {
// 		if index >= numLeafs/2 {
// 			proof = append(proof, node.Left.Hash)
// 			nodes = append(nodes, *node.Left)
// 			node = node.Right
// 			index = index - numLeafs/2
// 		} else {
// 			proof = append(proof, node.Right.Hash)
// 			nodes = append(nodes, *node.Right)
// 			node = node.Left
// 		}
// 		numLeafs = numLeafs / 2
// 	}

// 	return nodes
// }

func VerifyProof(root []byte, index int, hash []byte, proof [][]byte) (bool, error) {
	var err error
	for i := 0; i < len(proof); i++ {
		if index%2 == 0 {
			hash, err = GetHashFromBytes(append(hash, proof[i]...))
			if err != nil {
				return false, err
			}
		} else {
			hash, err = GetHashFromBytes(append(proof[i], hash...))
			if err != nil {
				return false, err
			}
		}
		index = index / 2
	}

	return bytes.Equal(hash, root), nil
}

func MarshalTree(tree *MerkleTree) ([]byte, error) {
	return json.Marshal(tree)
}

func UnmarshalTree(data []byte, tree *MerkleTree) error {
	return json.Unmarshal(data, tree)
}

func GetMerkleRoot(hashes [][]byte) ([]byte, error) {
	tree, err := NewMerkleTree(hashes)
	if err != nil {
		return nil, err
	}
	return tree.Root.Hash, nil
}

func NewMerkleTree(hashes [][]byte) (*MerkleTree, error) {

	var nodes = make([]*Node, 0, len(hashes))

	index := 0
	for _, hash := range hashes {
		nodes = append(nodes, &Node{Hash: hash})
		index++
	}

	level := nodes

	for len(level) > 1 {
		if len(level)%2 == 1 {
			index++
			level = append(level, level[len(level)-1])
		}

		newLevel := make([]*Node, 0, len(level)/2)
		for i := 0; i < len(level); i = i + 2 {
			newNode := &Node{Left: level[i], Right: level[i+1]}
			index++
			hash, err := GetHashFromBytes(append(newNode.Left.Hash, newNode.Right.Hash...))
			if err != nil {
				return nil, err
			}
			newNode.Hash = hash
			newLevel = append(newLevel, newNode)
		}

		level = newLevel
	}

	level[0].Num = &index
	return &MerkleTree{Root: level[0]}, nil
}

func GetHashFromBytes(data []byte) ([]byte, error) {
	hash := sha256.New()
	reader := bytes.NewReader(data)
	_, err := io.Copy(hash, reader)

	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func GetHashFromReader(reader io.Reader) ([]byte, error) {
	hash := sha256.New()

	_, err := io.Copy(hash, reader)

	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

type MerkleTree struct {
	Root *Node
}

type Node struct {
	Left  *Node
	Right *Node
	Hash  []byte
	Num   *int `json:"num,omitempty"`
}
