package crypto

import (
	"bytes"
	"testing"
)

func TestNewTPKE(t *testing.T) {
	tpkes := NewTPKE(4, 1)
	message := []byte("who are you")
	ct := tpkes[0].Enc(message)
	s0 := tpkes[0].DecShare(ct)
	s1 := tpkes[1].DecShare(ct)
	s2 := tpkes[2].DecShare(ct)
	s3 := tpkes[3].DecShare(ct)
	s := map[int]DecryptionShare{0: s0, 1: s1, 2: s2, 3: s3}
	msg := tpkes[1].Dec(s, ct)
	if !bytes.Equal(msg, message) {
		t.Errorf("decrypt error: %s", string(msg))
	}
}
func TestNewTPKE2(t *testing.T) {
	tpkes := NewTPKE(4, 1)
	message := []byte("who are you")
	ct := tpkes[0].Enc(message)
	s0 := tpkes[0].DecShare(ct)
	s1 := tpkes[1].DecShare(ct)
	s := map[int]DecryptionShare{0: s0, 1: s1}
	msg := tpkes[2].Dec(s, ct)
	if !bytes.Equal(msg, message) {
		t.Errorf("decrypt error: %s", string(msg))
	}
}
func TestNewTPKE3(t *testing.T) {
	tpkes := NewTPKE(4, 1)
	message := []byte("who are you")
	ct := tpkes[3].Enc(message)
	s0 := tpkes[0].DecShare(ct)
	s1 := tpkes[1].DecShare(ct)
	s := map[int]DecryptionShare{0: s0, 1: s1}
	msg := tpkes[2].Dec(s, ct)
	if !bytes.Equal(msg, message) {
		t.Errorf("decrypt error: %s", string(msg))
	}
}
func TestNewTPKE4(t *testing.T) {
	tpkes := NewTPKE(4, 1)
	message := []byte("who are you")
	ct := tpkes[3].Enc(message)
	s0 := tpkes[0].DecShare(ct)
	//s1 := tpkes[1].DecShare(ct)
	s := map[int]DecryptionShare{0: s0}
	msg := tpkes[2].Dec(s, ct)
	if bytes.Equal(msg, message) {
		t.Errorf("decrypt error: %s", string(msg))
	}
}
