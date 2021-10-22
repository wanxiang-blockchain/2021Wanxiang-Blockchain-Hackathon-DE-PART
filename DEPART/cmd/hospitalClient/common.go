package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"DEPART/goRecrypt/recrypt"

	"github.com/PlatONE_Network/PlatONE-SDK-Go/client"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/common"
	"github.com/PlatONE_Network/PlatONE-SDK-Go/platone/rpc"
)

const (
	IPFSUrl   = "127.0.0.1:5001"
	ServerUrl = "127.0.0.1:8080"
)

//发送GET请求
//url:请求地址
//response:请求返回的内容
func Get(url string) (response string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, error := client.Get(url)
	defer resp.Body.Close()
	if error != nil {
		panic(error)
	}

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	response = result.String()
	return
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func Post(url string, data interface{}) (content string) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}

func capsule2String(capsule recrypt.Capsule) string {
	capsuleAsBytes, _ := recrypt.EncodeCapsule(capsule)
	newCapsuleStr := hex.EncodeToString(capsuleAsBytes)
	return newCapsuleStr
}

func string2Capsule(capsuleStr string) *recrypt.Capsule {
	capsuleBytes, _ := hex.DecodeString(capsuleStr)
	capsule, _ := recrypt.DecodeCapsule(capsuleBytes)

	return &capsule
}

func InitContractClient() (common.TxParams, client.ContractClient) {
	txparam := common.TxParams{}
	keyfile := "./conf/keyfile.json"
	abiPath := "./conf/example/example.cpp.abi.json"
	PassPhrase := "0"
	vm := "wasm"
	url := client.URL{
		IP:      "127.0.0.1",
		RPCPort: 6791,
	}
	rpc, _ := rpc.DialContext(context.Background(), url.GetEndpoint())
	pc := client.Client{
		RpcClient:   rpc,
		Passphrase:  PassPhrase,
		KeyfilePath: keyfile,
		URL:         &url,
	}
	contract := client.ContractClient{
		Client:  &pc,
		AbiPath: abiPath,
		Vm:      vm,
	}
	return txparam, contract
}

func Hash(data string) []byte {
	h := sha256.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}
