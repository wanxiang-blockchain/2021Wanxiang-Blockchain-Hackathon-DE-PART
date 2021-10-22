package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/urfave/cli"
	"goRecrypt/curve"
	"goRecrypt/recrypt"
	"os"
)
var aPriKey ecdsa.PrivateKey

func main() {
	// Alice Generate Alice key-pair
	aPriKey, aPubKey, _ := curve.GenerateKeys()
	// Bob Generate Bob key-pair
	bPriKey, bPubKey, _ := curve.GenerateKeys()
	// plain text
	m := "Hello, Proxy Re-Encryption"
	fmt.Println("origin message:", m)
	// Alice encrypts to get cipherText and capsule
	cipherText, capsule, err := recrypt.Encrypt(m, aPubKey)
	if err != nil {
		fmt.Println(err)
	}
	capsuleAsBytes, err := recrypt.EncodeCapsule(*capsule)
	if err != nil {
		fmt.Println("encode error:", err)
	}
	capsuleTest, err := recrypt.DecodeCapsule(capsuleAsBytes)
	if err != nil {
		fmt.Println("decode error:", err)
	}
	fmt.Println("capsule before encode:", capsule)
	fmt.Println("capsule after decode:", capsuleTest)
	fmt.Println("ciphereText:", cipherText)
	// Test recreate aes key
	keyBytes, err := recrypt.RecreateAESKeyByMyPriKey(capsule, aPriKey)
	if err != nil {
		fmt.Println("Recreate key error:", err)
	}
	fmt.Println("recreate key:", hex.EncodeToString(keyBytes))

	// Alice generates re-encryption key
	rk, pubX, err := recrypt.ReKeyGen(aPriKey, bPubKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("rk:", rk)
	// Server executes re-encrypt
	newCapsule, err := recrypt.ReEncryption(rk, capsule)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Bob decrypts the cipherText
	plainText, err := recrypt.Decrypt(bPriKey, newCapsule, pubX, cipherText)
	if err != nil {
		fmt.Println(err)
	}

	plainTextByMyPri, err := recrypt.DecryptOnMyPriKey(aPriKey, capsule, cipherText)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("PlainText by my own private key:", string(plainTextByMyPri))
	// get plainText
	fmt.Println("plainText:", string(plainText))

	fileCapsule, err := recrypt.EncryptFile("a.txt", "a_encrypt.txt", aPubKey)
	if err != nil {
		fmt.Println("File Encrypt Error:", err)
	}
	fileNewCapsule, err := recrypt.ReEncryption(rk, fileCapsule)
	if err != nil {
		fmt.Println("ReEncryption Error:", err)
	}
	err = recrypt.DecryptFile("a_encrypt.txt", "a_decrypt.txt", bPriKey, fileNewCapsule, pubX)
	if err != nil {
		fmt.Println("Decrypt Error:", err)
	}

	app := cli.NewApp()
	app.Name = "sysconfig"
	app.Usage = "Setting basic configuration"
	app.Version = "0.0.1"

	app.Action = func(c *cli.Context) error {
		fmt.Println("Prepare applying basic configuration")
		a,b,_ := genRKAB("121212")
		fmt.Println(a)
		fmt.Println(b)
		return nil
	}

	app.Run(os.Args)
}

// 用户授权 b 能够访问数据，输入b 的公钥
func genRKAB(bPubKeyInput string)  (string, string, error) {

	fmt.Println("用户输入的公钥是:",bPubKeyInput)
	bPubKey, err := PublicKeyStrToKey(bPubKeyInput)
	if err != nil{
		return "","",err
	}

	rk, pubX, err := recrypt.ReKeyGen(&aPriKey, bPubKey)

	if err != nil {
		return "","",err
	}
	fmt.Println("产生的重加密密钥rk是:", rk)
	pubXString := PublicKeyToString(pubX)
	fmt.Println("产生的PubX 是:", pubXString)


	return rk.String(), pubXString,nil
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