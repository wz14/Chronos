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
	s := map[int][]byte{0: s0, 1: s1, 2: s2, 3: s3}
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
	s := map[int][]byte{0: s0, 1: s1}
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
	s := map[int][]byte{0: s0, 1: s1}
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
	s := map[int][]byte{0: s0}
	msg := tpkes[2].Dec(s, ct)
	if bytes.Equal(msg, message) {
		t.Errorf("decrypt error: %s", string(msg))
	}
}

func TestTPKE_Marshal(t *testing.T) {
	tpkes := NewTPKE(4, 1)
	m, err := tpkes[1].Marshal()
	if err != nil {
		t.Errorf("marshal fail %s", err.Error())
	}
	newtpke := TPKE{}
	err = newtpke.UnMarshal(m)
	if err != nil {
		t.Errorf("unmarshal fail %s", err.Error())
	}
	message := []byte("0xdeadbeef")
	ct := tpkes[3].Enc(message)
	s0 := tpkes[0].DecShare(ct)
	s1 := newtpke.DecShare(ct)
	s := map[int][]byte{0: s0, 1: s1}
	msg := tpkes[2].Dec(s, ct)
	if !bytes.Equal(msg, message) {
		t.Errorf("decrypt error %s", message)
	}
}

func TestTPKE_UnMarshal(t *testing.T) {
	tpkes := NewTPKE(4, 1)
	m, err := tpkes[0].Marshal()
	if err != nil {
		t.Errorf("marshal fail %s", err.Error())
	}
	new0tpke := TPKE{}
	err = new0tpke.UnMarshal(m)
	if err != nil {
		t.Errorf("unmarshal fail %s", err.Error())
	}
	m, err = tpkes[1].Marshal()
	if err != nil {
		t.Errorf("marshal fail %s", err.Error())
	}
	new1tpke := TPKE{}
	err = new1tpke.UnMarshal(m)
	if err != nil {
		t.Errorf("unmarshal fail %s", err.Error())
	}
	message := []byte("0xdeadbeef")
	ct := tpkes[3].Enc(message)
	s0 := new0tpke.DecShare(ct)
	s1 := new1tpke.DecShare(ct)
	s := map[int][]byte{0: s0, 1: s1}
	msg := tpkes[2].Dec(s, ct)
	if !bytes.Equal(msg, message) {
		t.Errorf("decrypt error %s", message)
	}
}
