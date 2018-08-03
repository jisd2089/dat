package security

import "drcs/common/cncrypt"

type DigitalEnvelope struct {
	// 使用公钥加密过的秘钥,hex16
	EncryptedKey string
	// 使用秘钥加密得到的数据,hex16
	Content string
}

func PackEnvelope(pubKey string, text string) *DigitalEnvelope {
	envelope := cncrypt.EncrpytEnvelopSafe(text, pubKey)
	return &DigitalEnvelope{envelope.Deskey, envelope.Ciphertext}
}

func UnpackEnvelope(digitalEnvelope *DigitalEnvelope) (string, error) {
	encryptedKey := digitalEnvelope.EncryptedKey
	content := digitalEnvelope.Content
	privateKey, err := GetPrivateKey()
	if err != nil {
		return "", err
	}
	envelope := &cncrypt.Envelope{encryptedKey, content}
	return envelope.DecryptEnvelopSafe(privateKey), nil
}
