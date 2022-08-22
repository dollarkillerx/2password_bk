package utils

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"
)

func GetPassword(password string, sale string) string {
	return Md5Encode(fmt.Sprintf("%s_%s", password, sale))
}

// 获取md5
func Md5Encode(str string) string {
	data := []byte(str)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// Base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64解码
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Base64URL编码
func Base64URLEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64URL解码
func Base64URLDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

/*
# RAS

公鑰密鑰的非對稱加密算法

公鑰加密只能通過私鑰解密
*/

// GenerateRSAKey RSA256 公钥密钥对生成
// @params: bits 密钥长度
// @returns: private 密钥
// @returns: public 公钥
func GenerateRSAKey(bits int) (priKey string, pubKey string, err error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priBlock := &pem.Block{
		Type:  "RAS PRIVATE KEY",
		Bytes: derStream,
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	return string(pem.EncodeToMemory(priBlock)), string(pem.EncodeToMemory(publicBlock)), nil
}

// RSAEncrypt Rsa256 加密
// @params: origData 原始数据
// @Params: pubKey 公钥
func RSAEncrypt(origData, pubKey interface{}) ([]byte, error) {
	var origDataBytes []byte
	switch r := origData.(type) {
	case string:
		origDataBytes = []byte(r)
	case []byte:
		origDataBytes = r
	}

	var pubKeyBytes []byte
	switch r := pubKey.(type) {
	case string:
		pubKeyBytes = []byte(r)
	case []byte:
		pubKeyBytes = r
	}

	//解密pem格式的公钥
	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origDataBytes)
}

// RSADecrypt Rsa256 解密
// @params: ciphertext 加密数据
// @Params: prvKey 私钥
func RSADecrypt(ciphertext, privateKey interface{}) ([]byte, error) {
	var ciphertextBytes []byte
	switch r := ciphertext.(type) {
	case string:
		ciphertextBytes = []byte(r)
	case []byte:
		ciphertextBytes = r
	}

	var privateKeyBytes []byte
	switch r := privateKey.(type) {
	case string:
		privateKeyBytes = []byte(r)
	case []byte:
		privateKeyBytes = r
	}

	//解密
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertextBytes)
}

// RSASign Rsa256 签名
// @params: origData 需要签名的数据
// @Params: prvKey 私钥
func RSASign(data, prvKey interface{}) ([]byte, error) {
	var dataBytes []byte
	switch r := data.(type) {
	case string:
		dataBytes = []byte(r)
	case []byte:
		dataBytes = r
	}

	var privateKeyBytes []byte
	switch r := prvKey.(type) {
	case string:
		privateKeyBytes = []byte(r)
	case []byte:
		privateKeyBytes = r
	}

	hash := sha256.New()
	hash.Write(dataBytes)
	sum := hash.Sum(nil)

	//解密私钥
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 签名
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, sum)
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
