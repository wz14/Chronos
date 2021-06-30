package crypto

import (
	"bytes"
	"testing"
)

//whole test
func TestErasure(t *testing.T) {
	s, err := Erasure([]byte("who are you"), 4, 1)
	if err != nil {
		t.Fatalf("erasure generate fail: %s", err.Error())
	}

	msg, err := Recover(s, 4, 1)
	if err != nil {
		t.Fatalf("erasure recover fail: %s", err.Error())
	}

	if !bytes.Equal([]byte("who are you"), msg) {
		t.Fatalf("recover bytes: %s", string(msg))
	}
}

//part test
func TestErasure2(t *testing.T) {
	s, err := Erasure([]byte("who are you"), 4, 1)
	if err != nil {
		t.Fatalf("erasure generate fail: %s", err.Error())
	}

	msg, err := Recover(s[:2], 4, 1)
	if err != nil {
		t.Fatalf("erasure recover fail: %s", err.Error())
	}

	if !bytes.Equal([]byte("who are you"), msg) {
		t.Fatalf("recover bytes: %s", string(msg))
	}
}

func TestErasure3(t *testing.T) {
	s, err := Erasure([]byte("who are you"), 4, 1)
	if err != nil {
		t.Fatalf("erasure generate fail: %s", err.Error())
	}

	msg, err := Recover(s[1:3], 4, 1)
	if err != nil {
		t.Fatalf("erasure recover fail: %s", err.Error())
	}

	if !bytes.Equal([]byte("who are you"), msg) {
		t.Fatalf("recover bytes: %s", string(msg))
	}
}

func TestErasure4(t *testing.T) {
	s, err := Erasure([]byte("who are you"), 4, 1)
	if err != nil {
		t.Fatalf("erasure generate fail: %s", err.Error())
	}

	msg, err := Recover([][]byte{s[0], s[3]}, 4, 1)
	if err != nil {
		t.Fatalf("erasure recover fail: %s", err.Error())
	}

	if !bytes.Equal([]byte("who are you"), msg) {
		t.Fatalf("recover bytes: %s", string(msg))
	}
}

func TestErasure5(t *testing.T) {
	s, err := Erasure([]byte("who are you"), 32, 8)
	if err != nil {
		t.Fatalf("erasure generate fail: %s", err.Error())
	}

	msg, err := Recover(s[4:13], 32, 8)
	if err != nil {
		t.Fatalf("erasure recover fail: %s", err.Error())
	}

	if !bytes.Equal([]byte("who are you"), msg) {
		t.Fatalf("recover bytes: %s", string(msg))
	}
}
