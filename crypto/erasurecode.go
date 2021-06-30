package crypto

import (
	"acc/pb"
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/klauspost/reedsolomon"
	"github.com/pkg/errors"
)

//Index:                0,
//Msglength:            0,
//Merkleroot:           nil,
//Merklepath:           nil,
//Merkleindex:          nil,
//ErasureCode:          nil,

// Erasure is t can't recover msg, but t+1 can recover
// return N list of msg
func Erasure(msg []byte, N, t int) ([][]byte, error) {
	length := len(msg)
	threshold := t + 1
	paddingmsg := append(msg, bytes.Repeat([]byte{0}, (length/threshold+1)*threshold-length)...)
	blocklength := len(paddingmsg) / threshold

	// fill data
	data := make([][]byte, N)
	for i := 0; i < N; i++ {
		data[i] = make([]byte, blocklength)
	}
	for i := 0; i < threshold; i++ {
		data[i] = paddingmsg[i*blocklength : (i+1)*blocklength]
		//fmt.Printf("%d: %d\n",i ,len(data[i]))
	}

	enc, err := reedsolomon.New(threshold, N-threshold)
	if err != nil {
		return nil, errors.Wrap(err, "create reed solomon fail")
	}

	err = enc.Encode(data)
	if err != nil {
		return nil, errors.Wrap(err, "encode fail")
	}

	// calculate merkle tree
	//merkletree.NewTree(data)
	tre, err := NewTree(data)
	if err != nil {
		return nil, errors.Wrap(err, "create merkle tree fail")
	}

	root := tre.GetRoot()

	result := [][]byte{}
	for i, share := range data {
		path, index := tre.GetProof(i)
		d := pb.ECMsg{
			Index:       uint64(i),
			Msglength:   uint64(length),
			Merkleroot:  root,
			Merklepath:  path,
			Merkleindex: index,
			ErasureCode: share,
		}
		byt, err := proto.Marshal(&d)
		if err != nil {
			return nil, errors.Wrapf(err, "marshal fail in %d", i)
		}
		result = append(result, byt)
	}

	return result, nil
}

// Recover s includes at least t+1 []byte to recover
// return msg
func Recover(s [][]byte, N, t int) ([]byte, error) {
	threshold := t + 1
	if len(s) < threshold {
		return nil, errors.Errorf("too few part, need t+1 %d", threshold)
	}

	data := make([][]byte, N)
	for j := 0; j < N; j++ {
		data[j] = nil
	}

	var length uint64

	for _, parts := range s {
		ecm := &pb.ECMsg{}
		err := proto.Unmarshal(parts, ecm)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal fail")
		}
		if !VerifyTree(ecm.Merkleroot, ecm.Merklepath,
			ecm.Merkleindex, ecm.ErasureCode) {
			return nil, errors.Wrap(err, "verify merkle tree fail")
		}
		data[ecm.Index] = ecm.ErasureCode
		length = ecm.Msglength
	}

	enc, err := reedsolomon.New(threshold, N-threshold)
	if err != nil {
		return nil, errors.Wrap(err, "create reed solomon fail")
	}

	err = enc.ReconstructData(data)
	if err != nil {
		return nil, errors.Wrap(err, "reconstruct data fail")
	}

	paddingmsg := bytes.Join(data[:threshold], []byte(""))

	return paddingmsg[:length], nil

}
