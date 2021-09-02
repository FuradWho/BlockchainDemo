package main

import (
	"SimpleBlockchainDemoV5/base58"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//1. 创建一个结构WalletKeyPair秘钥对，保存公钥和私钥
//2. 给这个结构提供一个方法GetAddress：私钥->公钥->地址

type WalletKeyPair struct {
	PrivateKey *ecdsa.PrivateKey

	//将公钥的X，Y进行字节流拼接后传输，这样在对端再进行切割还原，好处是可以方便后面的编码
	PublicKey []byte
}

func NewWalletKeyPair() *WalletKeyPair {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	publicKeyRaw := privateKey.PublicKey

	publicKey := append(publicKeyRaw.X.Bytes(), publicKeyRaw.Y.Bytes()...)
	return &WalletKeyPair{PrivateKey: privateKey, PublicKey: publicKey}
}

func (w *WalletKeyPair) GetAddress() string {

	publicHash := HashPubKey(w.PublicKey)
	version := 0x00
	//21字节的数据
	payload := append([]byte{byte(version)}, publicHash...)
	checksum := CheckSum(payload)
	//25字节
	payload = append(payload, checksum...)
	address := base58.Encode(payload)
	return address
}

func HashPubKey(pubKey []byte) []byte {
	hash := sha256.Sum256(pubKey)

	//创建一个hash160对象
	//向hash160中write数据
	//做哈希运算

	rip160Haher := ripemd160.New()
	_, err := rip160Haher.Write(hash[:])

	if err != nil {
		log.Panic(err)
	}

	//Sum函数 Sum参数append到一起返回，传入nil
	publicHash := rip160Haher.Sum(nil)
	return publicHash
}

func CheckSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	//4字节校验码
	checksum := second[0:4]
	return checksum
}
