package cncrypt

import (
	"errors"
	"strings"
)

const _ENCRYPT_TIMES = 1

//const _define MAX_STR_LEN = 64

var key = [16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}

func EXID_to_ETID(exid string) (string, error) {
	var output []byte
	if (len(exid) % 32) != 0 {
		return "", errors.New("input invalid!!")
	}
	id := HexBytes(exid)
	//id.SetString(exid, 16)
	//fmt.Println("exid: ", id.Text(16))
	var sm4 = new(SM4)
	sm4.EncryptSetKey(key[:])
	for i := 0; i < _ENCRYPT_TIMES; i++ {
		output = sm4.Encrypt(id)
	}
	if isUpper(exid) {
		return strings.ToUpper(HexString(output)), nil
	}
	return strings.ToLower(HexString(output)), nil
}
func ETID_to_EXID(etid string) (string, error) {
	var output []byte
	if (len(etid) % 32) != 0 {
		return "", errors.New("inpu invalid!!")
	}
	id := HexBytes(etid)
	//id.SetString(etid, 16)
	//fmt.Println("etid: ", id.Text(16))
	var sm4 = new(SM4)
	sm4.DecryptSetKey(key[:])
	for i := 0; i < _ENCRYPT_TIMES; i++ {
		output = sm4.Decrypt(id)
	}
	if isUpper(etid) {
		return strings.ToUpper(HexString(output)), nil
	}
	return strings.ToLower(HexString(output)), nil
}
func isUpper(exid string) bool {
	for i := 0; i < len(exid); i++ {

		if exid[i] >= 'A' && exid[i] <= 'F' {
			return true
		}

	}
	return false
}
