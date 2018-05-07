package security

import "drcs/dep/cncrypt"

// Signature 数字签名函数
func Signature(content []byte) (string, error) {
	priKey, err := GetPrivateKey()
	if err != nil {
		return "", err
	}
	return cncrypt.Sign(priKey, content), nil
}

// VerifySignature 校验签名，成功时返回true，否则返回false
func VerifySignature(pubKey string, content []byte, signature string) bool {
	result := cncrypt.VerifySign(pubKey, content, signature)
	if result == 1 {
		return true
	}
	return false
}
