package crypto

import (
	"bytes"
	"testing"
)

func TestNewTree(t *testing.T) {
	data := [][]byte{
		[]byte("hi"),
		[]byte("hello"),
		[]byte("fuck"),
		[]byte("who"),
	}
	tre, err := NewTree(data)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	b0, _ := tre.GetProof(0)
	b1, _ := tre.GetProof(1)

	if !bytes.Equal(b0[1], b1[1]) {
		t.Errorf("proof generate fail")
	}

	b2, _ := tre.GetProof(2)
	b3, _ := tre.GetProof(3)

	if !bytes.Equal(b2[1], b3[1]) {
		t.Errorf("proof generate fail")
	}
}

func TestVerifyTree(t *testing.T) {
	data := [][]byte{
		[]byte("hi"),
		[]byte("hello"),
		[]byte("fuck"),
		[]byte("who"),
	}
	tre, err := NewTree(data)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	root := tre.GetRoot()

	b0, indic0 := tre.GetProof(0)
	if !VerifyTree(root, b0, indic0, data[0]) {
		t.Errorf("verify fail")
	}
	b1, indic1 := tre.GetProof(1)
	if !VerifyTree(root, b1, indic1, data[1]) {
		t.Errorf("verify fail")
	}
	b2, indic2 := tre.GetProof(2)
	if !VerifyTree(root, b2, indic2, data[2]) {
		t.Errorf("verify fail")
	}
	b3, indic3 := tre.GetProof(3)
	if !VerifyTree(root, b3, indic3, data[3]) {
		t.Errorf("verify fail")
	}
}
