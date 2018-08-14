package transfer

/**
    Author: luzequan
    Created: 2018-08-01 10:13:09
*/
import (
	"drcs/core/interaction/response"
	"crypto/aes"
	"crypto/cipher"
	"bytes"
	"errors"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"fmt"
	"github.com/farmerx/gorsa"
	"encoding/base64"
)

type EncryptTransfer struct{}

func NewEncryptTransfer() Transfer {
	return &EncryptTransfer{}
}

// 封装Encryp服务
func (ft *EncryptTransfer) ExecuteMethod(req Request) Response {

	requestTxt := req.GetParameters()
	key := []byte(req.Param("encryptKey"))

	var (
		err        error
		body       []byte
		bodyByte   []byte
		bodyStr    string
		returnCode = "000000"
		retMsg     = "encryption or decryption success"
	)

	switch req.GetMethod() {
	case "AESENCRYPT":
		iv := []byte(req.Param("iv"))
		if req.Param("iv") == "" {
			iv = []byte(ivDefValue)
		}
		bodyByte, err = aesEncrypt(requestTxt, key, iv)
		bodyStr = base64.StdEncoding.EncodeToString(bodyByte)
	case "AESDECRYPT":
		ciphertext := req.GetPostData()
		requestTxt, err = base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return &response.DataResponse{
				StatusCode: 200,
				ReturnCode: "000005",
				ReturnMsg:  err.Error(),
			}
		}

		iv := []byte(req.Param("iv"))
		if req.Param("iv") == "" {
			iv = []byte(ivDefValue)
		}
		body, err = aesDecrypt(requestTxt, key, iv)

	case "RSAENCRYPT":
		bodyByte, err = rsaEncrypt(requestTxt, key)
		bodyStr = base64.StdEncoding.EncodeToString(bodyByte)
	case "RSADECRYPT":
		ciphertext := req.GetPostData()
		requestTxt, err = base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return &response.DataResponse{
				StatusCode: 200,
				ReturnCode: "000005",
				ReturnMsg:  err.Error(),
			}
		}
		body, err = rsaDecrypt(requestTxt, key)
	}

	if err != nil {
		returnCode = "000005"
		retMsg = err.Error()
	}

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: returnCode,
		ReturnMsg:  retMsg,
		Body:       body,
		BodyStr:    bodyStr,
	}
}

const (
	ivDefValue = "0102030405060708"
)

// AES 加解密 #####################################################################
func aesEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}
	plaintext = PKCS7Padding(plaintext, block.BlockSize())
	//iv := []byte(ivDefValue)
	blockMode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

func aesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}

	blockSize := block.BlockSize()

	if len(ciphertext) < blockSize {
		return nil, errors.New("ciphertext too short")
	}

	//iv := []byte(ivDefValue)
	fmt.Println("size", len(ciphertext)%blockSize)
	if len(ciphertext)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	blockModel := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plaintext, ciphertext)
	plaintext = PKCS7UnPadding(plaintext)

	return plaintext, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// RSA 加解密 #######################################################################
func rsaEncrypt(origData, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey) //将密钥解析成公钥实例
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData) //RSA算法加密
}

// 解密
func rsaDecrypt(ciphertext, privateKey []byte) ([]byte, error) {
	rsa := gorsa.RSA
	if err := rsa.SetPrivateKey(string(privateKey)); err != nil {
		return nil, err
	}

	return rsa.PriKeyDECRYPT(ciphertext)
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

//私钥
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`)

//公钥
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBEQqL3Hr7ud7MrEvuZMAVzl8C
jQwjK/sTx5UXDc+pUV+uIOhKA0wEG3Or+rH1wddITcW89Ti5zv+ypz1jlOtvS8GJ
+unjxxW7f4tLcmaUKWNxbhmgXZ6I05Dssa67oWhmPV/f5/L2Wgk9NFwbKYJWF7jP
UccC4+dC9f1FTroh5QIDAQAB
-----END PUBLIC KEY-----
`)

func (ft *EncryptTransfer) Close() {}
