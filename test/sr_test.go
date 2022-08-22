package test

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dollarkillerx/2password/internal/utils"
	"github.com/pkg/errors"
	"log"
	"testing"
)

var rk = `
-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgQD6mu0INiM1jMQEl0VloKme8wQ4MSsnlJR/b1eQ21cvczBaSvGB
ox0Ihc1lATSY4drcjtJLH/+s4bvIsE8Tp2DFK335Ex33BCm1uaS+vyEfw6ObqcAW
dHr7T2DzhUudzWf+g9nEQboxt4cgDN6UQh4pK/Ut6PQxI7+guzrkf1lKhwICEAEC
gYA18z0Wkz8g9eBXd5tXzLu6nmph83czcqdPe+H88WwPVrhd64SIwbOFPXwCn7h3
5CP+flki2xY8tlAqlE3CP0bBm8378FQgzYxSdW+sL2CIncLR0SRLJ9ZHMamA3/Z6
IEF9zpRvN0jhNVbSQuHbQQPaIapgsDhMb3O6gkWGqtjkQQJBAP8na2W7Ao0U9CwO
nX6ydbGLU6aGhmsIl88/jMO9dmw+6Hy8jVppgkk/xPcKj1/GXfEg1T0xgNZfQcZ9
I1fX+bcCQQD7b6UvpsNrmFNbPv9mGe+dw9rpE+B/taSyHWuLfNsw6luSHoL7Hl4Z
V3hmaHTxlSN6dTqpYOAd3y2mtNjlQ3WxAkBKm/3JhHoXYgrJtirtv5Yy6nFwPZS3
Qq87yfmsJUeubf9fha18+cxeZqBRvBrtdE11Dylrrx+cuVDH4Vugl0G1AkAlUD2B
Oq1XIyos6ItgcdJ0Q85BtNgFdJ8ofdYZUvNaDSjFJDUt86K9lyJtDLCQ0xYSzDnx
hUjv4CLEko8JF3IBAkAUKd2ggmiEzeMJNiHWlWR+pjVi1UedkjQOeDowJelOJTgy
jf1YhB/PxQ8jI4ZFSs0B8zoZOv95aYxITP/8ovw+
-----END RSA PRIVATE KEY-----
`

var pk = `
-----BEGIN RSA PUBLIC KEY-----
MIGIAoGBAPqa7Qg2IzWMxASXRWWgqZ7zBDgxKyeUlH9vV5DbVy9zMFpK8YGjHQiF
zWUBNJjh2tyO0ksf/6zhu8iwTxOnYMUrffkTHfcEKbW5pL6/IR/Do5upwBZ0evtP
YPOFS53NZ/6D2cRBujG3hyAM3pRCHikr9S3o9DEjv6C7OuR/WUqHAgIQAQ==
-----END RSA PUBLIC KEY-----
`

func TestRc(t *testing.T) {
	var sig = "U0bvkiOdzSjdQ2RcYgJB5T9eOuXKT2w2PXuPV6ZLz5WUre5Coa48M8KwX9biHhhJ+5GtEJAEUbwZJW/4k6j8L8SOsCOJBP4G/ZNaf4ZaVt1Xc4NZmsMjYKQeoX59IZKgR4rqpEbIxbajGqWy0btzbwob1kAc9P9b2jpVJFuv5bU="
	var data = "hello world"

	//sign, err := utils.RSASign(data, rk)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//fmt.Println(sig)
	//fmt.Println(utils.Base64Encode(sign))
	decode, err := utils.Base64Decode(sig)
	if err != nil {
		log.Fatalln(err)
	}
	err = RSASignVer(data, decode, pk)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("OK")
}

// RSASignVer Rsa256 验签
// @params: data 原始数据
// @params: signature 签名
// @params: publicKey 公钥
func RSASignVer(data, signature, publicKey interface{}) error {
	var dataBytes []byte
	switch r := data.(type) {
	case string:
		dataBytes = []byte(r)
	case []byte:
		dataBytes = r
	}

	var signatureBytes []byte
	switch r := signature.(type) {
	case string:
		signatureBytes = []byte(r)
	case []byte:
		signatureBytes = r
	}

	var publicKeyBytes []byte
	switch r := publicKey.(type) {
	case string:
		publicKeyBytes = []byte(r)
	case []byte:
		publicKeyBytes = r
	}

	hashed := sha256.Sum256(dataBytes)
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return err
	}
	//验证签名
	return rsa.VerifyPKCS1v15(pubInterface, crypto.SHA256, hashed[:], signatureBytes)
}
