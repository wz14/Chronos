package crypto

import "testing"

func TestNewCCconfigs(t *testing.T) {
	sigs := [][]byte{}
	configs, err := NewCCconfigs(2, 5)
	if err != nil {
		t.Fatal("config create fail")
	}
	message := []byte("balaabablablablaj")
	for i, c := range configs {
		sig, err := c.Sign(message)
		if err != nil {
			t.Fatalf("sign fail for %d", i)
		}
		sigs = append(sigs, sig)
	}
	ok, err := configs[0].Verify(sigs, message)
	if !ok {
		t.Errorf("verify fail with: %s", err.Error())
	}
}
