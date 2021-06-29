package crypto

// Note: code copy from
// file1: https://github.com/DE-labtory/cleisthenes/tpke/threshold_encrpytion.go
// file2: https://github.com/DE-labtory/cleisthenes/tpke/threshold_encryption_test.go
// with Apache 2.0 license
//
// Copyright 2019 DE-labtory
//
// Modify: combine file1 and file2, add keyGen algorithm,
//		   rename interface name to avoid repetition with
//		   this project package name

import (
	tpk "github.com/WangZhuo2000/tpke"
	"strconv"
)

type SecretKey [32]byte
type PublicKey []byte
type DecryptionShare [96]byte
type CipherText []byte

type Config struct {
	threshold   int
	participant int
}

type DefaultTpke struct {
	threshold    int
	publicKey    *tpk.PublicKey
	publicKeySet *tpk.PublicKeySet
	secretKey    *tpk.SecretKeyShare
	decShares    map[string]*tpk.DecryptionShare
}

func NewDefaultTpke(th int, skStr SecretKey, pksStr PublicKey) (*DefaultTpke, error) {
	sk := tpk.NewSecretKeyFromBytes(skStr)
	sks := tpk.NewSecretKeyShare(sk)

	pks, err := tpk.NewPublicKeySetFromBytes(pksStr)
	if err != nil {
		return nil, err
	}

	return &DefaultTpke{
		threshold:    th,
		publicKeySet: pks,
		publicKey:    pks.PublicKey(),
		secretKey:    sks,
		decShares:    make(map[string]*tpk.DecryptionShare),
	}, nil
}

func (t *DefaultTpke) AcceptDecShare(id string, decShare DecryptionShare) {
	ds := tpk.NewDecryptionShareFromBytes(decShare)
	t.decShares[id] = ds
}

func (t *DefaultTpke) ClearDecShare() {
	t.decShares = make(map[string]*tpk.DecryptionShare)
}

// Encrypt encrypts some byte array message.
func (t *DefaultTpke) Encrypt(msg []byte) ([]byte, error) {
	encrypted, err := t.publicKey.Encrypt(msg)
	if err != nil {
		return nil, err
	}
	return encrypted.Serialize(), nil
}

// DecShare makes decryption share using each secret key.
func (t *DefaultTpke) DecShare(ctb CipherText) DecryptionShare {
	ct := tpk.NewCipherTextFromBytes(ctb)
	ds := t.secretKey.DecryptShare(ct)
	return ds.Serialize()
}

// Decrypt collects decryption share, and combine it for decryption.
func (t *DefaultTpke) Decrypt(decShares map[string]DecryptionShare, ctBytes []byte) ([]byte, error) {
	ct := tpk.NewCipherTextFromBytes(ctBytes)
	ds := make(map[string]*tpk.DecryptionShare)
	for id, decShare := range decShares {
		ds[id] = tpk.NewDecryptionShareFromBytes(decShare)
	}
	return t.publicKeySet.DecryptUsingStringMap(ds, ct)
}

type TPKE struct {
	sks  *tpk.SecretKeySet
	pks  *tpk.PublicKeySet
	id   int
	tpke *DefaultTpke
}

func NewTPKE(N int, t int) []*TPKE {
	tpkes := []*TPKE{}
	c := &Config{
		threshold:   t,
		participant: N,
	}
	secretKeySet := tpk.RandomSecretKeySet(c.threshold)
	publicKeySet := secretKeySet.PublicKeySet()
	for i := 0; i < N; i++ {
		tpke, _ := NewDefaultTpke(N, secretKeySet.KeyShareUsingString(strconv.Itoa(i)).Serialize(),
			publicKeySet.Serialize())
		tpkes = append(tpkes, &TPKE{
			sks:  secretKeySet,
			pks:  publicKeySet,
			id:   i,
			tpke: tpke,
		})
	}
	return tpkes
}

func (t *TPKE) Enc(msg []byte) []byte {
	decshare, _ := t.tpke.Encrypt(msg)
	return decshare
}

func (t *TPKE) DecShare(ct []byte) []byte {
	ds := t.tpke.DecShare(ct)
	return ds[:]
}

func (t *TPKE) Dec(m map[int][]byte, ct []byte) []byte {
	mm := map[string]DecryptionShare{}
	for k, v := range m {
		var arr DecryptionShare
		copy(arr[:], v)
		mm[strconv.Itoa(k)] = arr
	}
	msg, _ := t.tpke.Decrypt(mm, ct)
	return msg
}
