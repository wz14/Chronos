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

func TestNewCCconfigs2(t *testing.T) {
	sigs := [][]byte{}
	configs, err := NewCCconfigs(2, 5)
	if err != nil {
		t.Fatal("config create fail")
	}
	message := []byte("balaabablablablaj")
	for i := 0; i < 3; i++ {
		sig, err := configs[i].Sign(message)
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

func TestNewCCconfigs3(t *testing.T) {
	sigs := [][]byte{}
	configs, err := NewCCconfigs(2, 5)
	if err != nil {
		t.Fatal("config create fail")
	}
	message := []byte("balaabablablablaj")
	for i := 0; i < 2; i++ {
		sig, err := configs[i].Sign(message)
		if err != nil {
			t.Fatalf("sign fail for %d", i)
		}
		sigs = append(sigs, sig)
	}
	ok, err := configs[0].Verify(sigs, message)
	if ok {
		t.Errorf("verify shouldn't be ok: %s", err.Error())
	}
}

func TestCCconfig_Marshal(t *testing.T) {
	sigs := [][]byte{}
	configs, err := NewCCconfigs(2, 5)
	if err != nil {
		t.Fatal("config create fail")
	}

	configByts, err := configs[0].Marshal()
	if err != nil {
		t.Fatalf("marshall fail: %s", err.Error())
	}
	new0Config := CCconfig{}
	err = new0Config.UnMarshal(configByts)
	if err != nil {
		t.Fatalf("unmarshall fail: %s", err.Error())
	}

	configByts, err = configs[1].Marshal()
	if err != nil {
		t.Fatalf("marshall fail: %s", err.Error())
	}
	new1Config := CCconfig{}
	err = new1Config.UnMarshal(configByts)
	if err != nil {
		t.Fatalf("unmarshall fail: %s", err.Error())
	}

	message := []byte("balaabablablablaj")

	sig, err := new0Config.Sign(message)
	if err != nil {
		t.Fatal("sign fail for 0")
	}
	sigs = append(sigs, sig)

	sig, err = new1Config.Sign(message)
	if err != nil {
		t.Fatal("sign fail for 1")
	}
	sigs = append(sigs, sig)

	sig, err = configs[3].Sign(message)
	if err != nil {
		t.Fatal("sign fail for 1")
	}
	sigs = append(sigs, sig)

	ok, err := configs[0].Verify(sigs, message)
	if !ok {
		t.Errorf("verify fail with: %s", err.Error())
	}
}
