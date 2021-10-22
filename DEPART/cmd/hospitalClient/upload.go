package main

import (
	"DEPART/goRecrypt/utils"
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"

	"gopkg.in/urfave/cli.v1"

	"DEPART/goRecrypt/recrypt"

	shell "github.com/ipfs/go-ipfs-api"
)

func handleUpload(c *cli.Context) {
	data := c.String("data")
	patientID := c.String("patientId")
	patientPkStr := c.String("patientPK")
	patientPK, err := utils.PublicKeyStrToKey(patientPkStr)
	if err != nil {
		panic(err)
	}

	CData, capsule, err := encryptData(data, patientPK)
	if err != nil {
		return
	}

	CDataUri, err := uploadToIPFS(CData)
	if err != nil {
		return
	}

	uploadToServer(patientID, CDataUri, capsule2String(*capsule))
	uploadToPlatone(data)
}

func encryptData(data string, patientPK *ecdsa.PublicKey) (CData []byte, capsule *recrypt.Capsule, err error) {
	CData, capsule, err = recrypt.Encrypt(data, patientPK)
	return
}

// C_data：ciphertext of data
func uploadToIPFS(CData []byte) (CDataUri string, err error) {
	sh := shell.NewShell(IPFSUrl)
	CDataUri, err = sh.Add(bytes.NewReader(CData))
	if err != nil {
		return "", err
	}
	return
}

func uploadToServer(patientID string, CDataUri string, capsule string) {
	Post(ServerUrl+"/uploadID", struct {
		ID      string `json:"id"`
		DataUri string `json:"dataUri"`
		Capsule string `json:"capsule"`
	}{
		ID:      patientID,
		DataUri: CDataUri,
		Capsule: capsule,
	})
}

func uploadToPlatone(data string) {
	dataHash := Hash(data)
	txparam, contract := InitContractClient()
	funcname := "setEvidence"
	funcparam := []string{hex.EncodeToString(dataHash), hex.EncodeToString(dataHash)}
	cns := "wxbc1"
	_, err := contract.Execute(context.Background(), txparam, funcname, funcparam, cns)
	if err != nil {
		panic(err)
	}

	return
}
