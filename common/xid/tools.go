package xid

import (
	"bytes"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"drcs/common/cncrypt"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"sort"
	"strings"
)

/*
以下为加密算法
//sign_type
SIGN_HMAC_SHA1 = 1
SIGN_HMAC_SHA256 = 2
SIGN_HMAC_MD5 = 3
SIGN_HMAC_SM3 = 4
//encrypt_type
ALG_3DES_ECB_PKCS5Padding = 1
ALG_SM4_ECB_PKCS5Padding = 2
//公安三所采用base64进行[]byte和string转换
*/
func sign_SHA256(key []byte, message string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func sign_MD5(key []byte, message string) string {
	mac := hmac.New(md5.New, key)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

//SM4_ECB_PKCS5Padding
func encryptSM4(key []byte, message string) (string, error) {
	//key should always be 24 bytes
	ckey := make([]byte, 24)
	copy(ckey, key)
	//blocksize for sm4 = 16, block for des = 8, be careful!!!
	blocksize := 16
	data := pkcs5padding([]byte(message), blocksize)

	var sm4 = new(cncrypt.SM4)
	sm4.EncryptSetKey(key[:])
	data = sm4.Encrypt(data[:])
	return base64.StdEncoding.EncodeToString(data), nil
}

//SM4_ECB_PKCS5Padding
func decryptSM4(key []byte, message string) (string, error) {
	//key should always be 24 bytes
	ckey := make([]byte, 24)
	copy(ckey, key)
	var sm4 = new(cncrypt.SM4)
	sm4.DecryptSetKey(key[:])
	data, _ := base64.StdEncoding.DecodeString(message)
	data = sm4.Decrypt(data[:])
	data = pkcs5unpadding(data)
	return string(data), nil
}

//ECBPKCS5Padding加密
func encrypt3DES(key []byte, message string) (string, error) {
	//key should always be 24 bytes
	ckey := make([]byte, 24)
	copy(ckey, key)
	block, err := des.NewTripleDESCipher(ckey)
	if err != nil {
		return "", err
	}
	blocksize := block.BlockSize()
	data := pkcs5padding([]byte(message), blocksize)
	block.Encrypt(data, data)

	return base64.StdEncoding.EncodeToString(data), nil
}

//ECBPKCS5Padding解密
func decrypt3DES(key []byte, message string) (string, error) {
	ckey := make([]byte, 24)
	copy(ckey, key)
	block, err := des.NewTripleDESCipher(ckey)
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}
	block.Decrypt(data, data)
	msg := pkcs5unpadding(data)
	return string(msg), nil
}

//ECB PKCS5Padding
func pkcs5padding(ciphertext []byte, blocksize int) []byte {
	paddingsize := blocksize - len(ciphertext)%blocksize
	padding := bytes.Repeat([]byte{byte(paddingsize)}, paddingsize)
	return append(ciphertext, padding...)
}

//ECB PKCS5Unpadding
func pkcs5unpadding(data []byte) []byte {
	length := len(data)
	unpaddingsize := int(data[length-1])
	return data[:(length - unpaddingsize)]
}

//消息摘要算法密钥生成
func genSignKey(sign_factor, appkey string) ([]byte, error) {
	if len(sign_factor) != 16 {
		return nil, errors.New("sign_factor should be 8 bytes in hex format")
	}
	if len(appkey) != 32 {
		return nil, errors.New("appkey should be 16 bytes in hex format")
	}
	return genkey(sign_factor, appkey), nil
}

//消息加密解密算法密钥生成
func genEncryptKey(encrypt_factor, appkey string) ([]byte, error) {
	if len(encrypt_factor) != 16 {
		return nil, errors.New("encrypt_factor should be 8 bytes in hex format")
	}
	if len(appkey) != 32 {
		return nil, errors.New("appkey should be 16 bytes in hex format")
	}
	return genkey(encrypt_factor, appkey), nil
}

//factor and appkey are hex-format strings
func genkey(factor, appkey string) []byte {
	factor_bytes, _ := hex.DecodeString(factor)
	data := [16]byte{}
	for i := 0; i < 8; i++ {
		data[i] = factor_bytes[i]
		data[i+8] = ^factor_bytes[i]
	}
	key, _ := hex.DecodeString(appkey)
	var sm4 = new(cncrypt.SM4)
	sm4.EncryptSetKey(key[:])

	return sm4.Encrypt(data[:])
}

func genrandomfactor() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return hex.EncodeToString(b)
}

func sortReqGenXid(genxid ReqGenXid) string {
	rjson, _ := json.Marshal(genxid)
	var mmap map[string]interface{}
	json.Unmarshal(rjson, &mmap)
	res := sortJsonString(mmap)
	return res
}

func sortReqConvXid(convid ReqConvXid) string {
	rjson, _ := json.Marshal(convid)
	var mmap map[string]interface{}
	json.Unmarshal(rjson, &mmap)
	res := sortJsonString(mmap)
	return res
}

func sortJsonString(data map[string]interface{}) string {
	mmap := data
	keys := make([]string, 0)
	for k, _ := range mmap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var res string
	for _, key := range keys {
		if key == "security_factor" {
			res += key + "=" + sortJsonString(mmap[key].(map[string]interface{})) + "&"
			continue
		}
		res += key + "=" + mmap[key].(string) + "&"
	}
	return strings.TrimRight(res, "&")
}

//generate uuid for biz_sequence_id
func genuuid() string {
	b := make([]byte, 48)
	io.ReadFull(rand.Reader, b)
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
	return hex.EncodeToString(h.Sum(nil))
}
