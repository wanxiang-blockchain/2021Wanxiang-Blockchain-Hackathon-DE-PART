package rpc

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"preServer/db"
	"preServer/recrypt"
	"preServer/utils"
)

func uploadId(rec *utils.Record) error {
	db.Put(rec)
	return nil
}

func reEncryption(key *utils.ReKey) (*utils.CapsuleResult, error) {

	rk := big.NewInt(0)
	_, ok := rk.SetString(key.Rk, 10)
	if !ok {
		return nil, fmt.Errorf("rk err")
	}
	rec := db.Get(key.ID)
	if rec == nil {
		return nil, fmt.Errorf("ID not exist")
	}
	capsuleBytes, err := hex.DecodeString(rec.Capsule)
	if err != nil {
		return nil, fmt.Errorf("decode capsule error")
	}
	capsule, err := recrypt.DecodeCapsule(capsuleBytes)
	if err != nil {
		return nil, fmt.Errorf("decode capsule error")
	}

	newCapsule, err := recrypt.ReEncryption(rk, &capsule)
	if err != nil {
		return nil, fmt.Errorf("reEncryption error")
	}

	capsuleAsBytes, _ := recrypt.EncodeCapsule(*newCapsule)

	newCapsuleStr := hex.EncodeToString(capsuleAsBytes)

	return &utils.CapsuleResult{
		NewCapsule: newCapsuleStr,
		DataUri:    rec.DataUri,
	}, nil
}
