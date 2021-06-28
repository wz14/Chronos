package crypto

import (
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/bls"
	"go.dedis.ch/kyber/v3/sign/tbls"
)

type PartPriKey = share.PriShare
type PubKey = share.PubPoly

func NewCCconfigs(t int, n int) ([]*CCconfig, error) {
	pskGroup, pk, err := internalKeyGen(t, n)
	if err != nil {
		return nil, errors.Wrap(err, "keygen fail")
	}
	cccs := []*CCconfig{}
	for _, psk := range pskGroup {
		cccs = append(cccs, &CCconfig{
			T:   t + 1, // t can't recover, and t+1 can recover
			N:   n,
			psk: psk,
			pk:  pk,
		})
	}
	return cccs, nil
}

// common coin config
type CCconfig struct {
	T   int
	N   int
	psk *PartPriKey
	pk  *PubKey
}

func (c *CCconfig) Sign(msg []byte) ([]byte, error) {
	return internalSign(c.psk, msg)
}

func (c *CCconfig) Verify(sigs [][]byte, msg []byte) (bool, error) {
	return internalVerfiy(c.pk, msg, sigs, c.T, c.N)
}

func (c *CCconfig) Combine(sigs [][]byte, msg []byte) ([]byte, error) {
	suit := pairing.NewSuiteBn256()
	sig, err := tbls.Recover(suit, c.pk, msg, sigs, c.T, c.N)
	if err != nil {
		return nil, errors.Wrap(err, "sig recover fail")
	}
	return sig, nil
}

// sklist, pk, error
func internalKeyGen(t int, n int) ([]*PartPriKey, *PubKey, error) {
	suit := pairing.NewSuiteBn256()
	random := suit.RandomStream()
	x := suit.G1().Scalar().Pick(random)

	// priploy
	priploy := share.NewPriPoly(suit.G2(), t, x, suit.RandomStream())
	// n points in ploy
	npoints := priploy.Shares(n)
	//pub ploy
	pubploy := priploy.Commit(suit.G2().Point().Base())
	return npoints, pubploy, nil
}

func internalSign(sk *PartPriKey, msg []byte) ([]byte, error) {
	suit := pairing.NewSuiteBn256()
	partSig, err := tbls.Sign(suit, sk, msg)
	if err != nil {
		return nil, err
	}
	return partSig, nil
}

func internalVerfiy(pk *PubKey, msg []byte, sigs [][]byte, t, n int) (bool, error) {
	suit := pairing.NewSuiteBn256()
	sig, err := tbls.Recover(suit, pk, msg, sigs, t, n)
	if err != nil {
		return false, errors.Wrap(err, "sig recover fail")
	}
	err = bls.Verify(suit, pk.Commit(), msg, sig)
	if err != nil {
		return false, errors.Wrap(err, "sig verify fail")
	}
	return true, nil
}
