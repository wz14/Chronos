package crypto

import (
	"bytes"
	"crypto/sha256"
	m "github.com/cbergoon/merkletree"
)

// implement of m.Content
type implContent struct {
	x []byte
}

func buildImplContent(x []byte) *implContent {
	return &implContent{x: x}
}

func (i *implContent) CalculateHash() ([]byte, error) {
	hash := Hash(i.x)
	return hash, nil
}

func (i *implContent) Equals(other m.Content) (bool, error) {
	hash1, _ := other.CalculateHash()
	hash2, _ := i.CalculateHash()
	if bytes.Equal(hash1, hash2) {
		return true, nil
	} else {
		return false, nil
	}
}

type Tree struct {
	mktree   *m.MerkleTree
	contents []m.Content
}

func NewTree(data [][]byte) (*Tree, error) {
	contents := []m.Content{}
	for _, d := range data {
		c := buildImplContent(d)
		contents = append(contents, c)
	}
	mk, err := m.NewTree(contents)
	if err != nil {
		return nil, err
	}
	return &Tree{
		mktree:   mk,
		contents: contents,
	}, nil
}

func (t *Tree) GetRoot() []byte {
	return t.mktree.MerkleRoot()
}

func (t *Tree) GetProof(id int) ([][]byte, []int64) {
	path, indicator, _ := t.mktree.GetMerklePath(t.contents[id])
	return path, indicator
}

func VerifyTree(root []byte, proof [][]byte, indicator []int64, msg []byte) bool {
	if len(proof) != len(indicator) {
		return false
	}
	itHash, _ := (&implContent{x: msg}).CalculateHash()
	for i, p := range proof {
		s := sha256.New()
		if indicator[i] == 1 {
			s.Write(append(itHash, p...))
		} else if indicator[i] == 0 {
			s.Write(append(p, itHash...))
		} else {
			return false
		}
		itHash = s.Sum(nil)
	}
	return bytes.Equal(itHash, root)
}
