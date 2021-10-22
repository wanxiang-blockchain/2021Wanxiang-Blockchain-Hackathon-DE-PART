package main

import (
	"DEPART/goRecrypt/utils"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/urfave/cli.v1"

	"DEPART/goRecrypt/recrypt"

	shell "github.com/ipfs/go-ipfs-api"
)

func handleDownload(c *cli.Context) (data string) {
	patientID := c.String("patientId")
	rk_A_B := c.String("rk")
	myPrivateKeyStr := c.String("myPrivateKey")
	myPrivateKey, err := utils.PrivateKeyStrToKey(myPrivateKeyStr)
	if err != nil {
		panic(err)
	}

	pubXStr := c.String("pubX")
	pubX, err := utils.PublicKeyStrToKey(pubXStr)
	if err != nil {
		panic(err)
	}

	capsuleStr, CDataUri := getFromServer(rk_A_B, patientID)
	CData := getFromIPFS(CDataUri)

	dataBytes, err := recrypt.Decrypt(myPrivateKey, string2Capsule(capsuleStr), pubX, CData)
	if err != nil {
		panic(err)
	}
	data = string(dataBytes)
	checkDataFromPlatone(data)
	return data
}

func getFromServer(rk_A_B string, patientID string) (capsule string, CDataUri string) {
	response := Post(ServerUrl+"/reEncryption", struct {
		ID string `json:"id"`
		Rk string `json:"rk"`
	}{
		ID: patientID,
		Rk: rk_A_B,
	})

	result := struct {
		NewCapsule string `json:"newCapsule"`
		DataUri    string `json:"dataUri"`
	}{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		panic(err)
	}

	return result.NewCapsule, result.DataUri
}

func getFromIPFS(CDataUri string) (CData []byte) {
	sh := shell.NewShell(IPFSUrl)
	err := sh.Get(CDataUri, "./")
	if err != nil {
		panic(err)
	}

	f, _ := os.Open(CDataUri)
	content, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return content
}

func checkDataFromPlatone(data string) {
	dataHash := Hash(data)

	txparam, contract := InitContractClient()
	funcname := "getEvidence"
	funcparam := []string{hex.EncodeToString(dataHash)}
	cns := "wxbc1"
	res, err := contract.Execute(context.Background(), txparam, funcname, funcparam, cns)
	if err != nil {
		panic(err)
	}

	hash := fmt.Sprintf("%v", res)
	if hex.EncodeToString(dataHash) != hash {

	}
}
