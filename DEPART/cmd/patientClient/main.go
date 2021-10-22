package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"goRecrypt/curve"
	"goRecrypt/recrypt"
	"os"
)

var (
	DefaultClient GenRkClient
)

func init() {
	DefaultClient = newGenRk()
}

type GenRkClient struct {
	aPriKey *ecdsa.PrivateKey
}

func newGenRk() GenRkClient {
	aPriKey, _ := GenAliceKeyPair()
	return GenRkClient{aPriKey: aPriKey}
}

func main() {
	_, bPubKey := GenBobKeyPair()
	var pubKey string
	app := cli.NewApp()
	app.Name = "DepartForAlice"
	app.Usage = "Setting basic configuration"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "bobPubKey, pk",
			Usage:       "input bob's publick key",
			Value:       PublicKeyToString(bPubKey),
			Destination: &pubKey,
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("Bob 的公钥是：", pubKey)
		rk, pubX, err := DefaultClient.GenReEncryptionKey(pubKey)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("产生的重加密密钥rk是:", rk)
		fmt.Println("产生的PubX 是:", pubX)
		return nil
	}

	app.Run(os.Args)
}

func GenAliceKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	// Alice Generate Alice key-pair
	aPriKey, aPubKey, _ := curve.GenerateKeys()
	return aPriKey, aPubKey
}

func GenBobKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	// Bob Generate Bob key-pair
	bPriKey, bPubKey, _ := curve.GenerateKeys()
	return bPriKey, bPubKey
}

// 用户授权 b 能够访问数据，输入b 的公钥
func (sk GenRkClient) GenReEncryptionKey(bPubKeyInput string) (string, string, error) {
	bPubKey, err := PublicKeyStrToKey(bPubKeyInput)
	if bPubKey.X == nil || bPubKey.Y == nil {
		fmt.Errorf("publick key struct is error, please input again")
		return "", "", errors.New("publick key struct is error")
	}
	if err != nil {
		return "", "", err
	}
	rk, pubX, err := recrypt.ReKeyGen(sk.aPriKey, bPubKey)
	if err != nil {
		return "", "", err
	}
	pubXString := PublicKeyToString(pubX)
	return rk.String(), pubXString, nil
}

// convert public key to string
func PublicKeyToString(publicKey *ecdsa.PublicKey) string {
	pubKeyBytes := curve.PointToBytes(publicKey)
	return hex.EncodeToString(pubKeyBytes)
}

// convert public key string to key
func PublicKeyStrToKey(pubKey string) (*ecdsa.PublicKey, error) {
	pubKeyAsBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(curve.CURVE, pubKeyAsBytes)
	key := &ecdsa.PublicKey{
		Curve: curve.CURVE,
		X:     x,
		Y:     y,
	}
	return key, nil
}
