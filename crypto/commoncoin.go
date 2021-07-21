package crypto

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/bls"
	"go.dedis.ch/kyber/v3/sign/tbls"
	"strconv"
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

func (c *CCconfig) Marshal() ([]string, error) {
	result := make([]string, 6)
	result[0] = strconv.Itoa(c.T)
	result[1] = strconv.Itoa(c.N)
	result[2] = strconv.Itoa(c.psk.I)
	// marshal psk
	byts, err := c.psk.V.MarshalBinary()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to marshal CCconfig.psk.V")
	}
	result[3] = base64.StdEncoding.EncodeToString(byts)

	base, committs := c.pk.Info()
	// marshal base
	byts, err = base.MarshalBinary()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to marshal CCconfig.pk.base")
	}
	result[4] = base64.StdEncoding.EncodeToString(byts)
	// marshal committs
	result[5] = strconv.Itoa(len(committs))
	for i, commit := range committs {
		byts, err = commit.MarshalBinary()
		if err != nil {
			return nil, errors.Wrapf(err, "fail to marshal CCconfig.pk.commit[%d]", i)
		}
		result = append(result, base64.StdEncoding.EncodeToString(byts))
	}
	return result, nil
}

func (c *CCconfig) UnMarshal(s []string) error {
	suit := pairing.NewSuiteBn256()
	var err error
	c.T, err = strconv.Atoi(s[0])
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal T")
	}
	c.N, err = strconv.Atoi(s[1])
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal N")
	}
	// unmarshal psk
	i, err := strconv.Atoi(s[2])
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal psk.i")
	}
	vbytes, err := base64.StdEncoding.DecodeString(s[3])
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal psk.v")
	}
	v := suit.G1().Scalar()
	err = v.UnmarshalBinary(vbytes)
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal psk.v")
	}
	c.psk = &PartPriKey{
		I: i,
		V: v,
	}
	// ummarshal pk
	baseByts, err := base64.StdEncoding.DecodeString(s[4])
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal pk.base")
	}
	base := suit.G2().Point()
	err = base.UnmarshalBinary(baseByts)
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal pk.base")
	}
	committsLen, err := strconv.Atoi(s[5])
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal pk.committsLen")
	}
	if committsLen+6 != len(s) {
		return errors.Errorf("pk.committsLen = %d; len(s) = %d", committsLen, len(s))
	}
	committs := make([]kyber.Point, committsLen)
	for i := 0; i < committsLen; i++ {
		commit := suit.G2().Point()
		commitByts, err := base64.StdEncoding.DecodeString(s[6+i])
		if err != nil {
			return errors.Wrapf(err, "fail to unmarshal pk.committs[%d]", i)
		}
		err = commit.UnmarshalBinary(commitByts)
		if err != nil {
			return errors.Wrapf(err, "fail to unmarshal pk.committs[%d]", i)
		}
		committs[i] = commit
	}
	c.pk = share.NewPubPoly(suit.G2(), base, committs)
	return nil
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
