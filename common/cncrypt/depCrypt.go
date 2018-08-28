package cncrypt

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//var _depCrypt *DepCrypt

//type DepCrypt struct {
//	SM2_ *SM2
//	SM4_ *SM4
//}

//func GetDepCrypt() *DepCrypt {
//	if _depCrypt != nil {
//		return _depCrypt
//	} else {
//		_depCrypt = &DepCrypt{}
//		_depCrypt.SM2_ = new(SM2)
//		_depCrypt.SM4_ = new(SM4)
//		_depCrypt.SM2_.Init()
//	}

//	return _depCrypt
//}

var key_map map[string][]byte
var pubkey_map map[string][2][]byte
var sm2 *SM2
var m *sync.RWMutex

var key_sm4 = [16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}

type Envelope struct {
	Deskey     string
	Ciphertext string
}

func clearMap() {
	for k, _ := range key_map {
		delete(key_map, k)
	}
	for k, _ := range pubkey_map {
		delete(pubkey_map, k)
	}
}

func startTimer(f func()) {
	go func() {
		for {

			now := time.Now()
			next := now.Add(time.Hour * 24 * 7)
			next = time.Date(next.Year(), next.Month(), next.Day(), 2, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			f()
		}
	}()
}
func Init(privateKey string) {
	key_map = make(map[string][]byte)
	pubkey_map = make(map[string][2][]byte)
	sm2 = new(SM2)

	sm2.Init()
	m = new(sync.RWMutex)
	if sm2.SignCtx == nil {
		sm2.SignCtx = &sm2SignContext{}
		sm2.SignCtx.Init = 0
	}
	signInitContext(sm2.SM2Params, sm2.SignCtx, privateKey)
	//fmt.Println("cnrypt manul init")
}

//func init() {
//	key_map = make(map[string][]byte)
//	pubkey_map = make(map[string][2][]byte)
//	sm2 = new(SM2)
//	sm2.Init()
//	m = new(sync.RWMutex)
//	go startTimer(clearMap)
//	fmt.Println("cnrypt  init")
//}

func (envelope *Envelope) setData(ciphertext, deskey string) {

	envelope.Deskey = deskey
	envelope.Ciphertext = ciphertext

}
func genKey() []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	key := make([]byte, 16)
	for i := 0; i < 16; i++ {
		key[i] = (byte)(r.Intn(255))
	}
	return key[:]
}
func EncrpytEnvelop(text, pubkey string) *Envelope {
	var sm4 = new(SM4)
	var key []byte
	//var ok bool
	if keys, ok := pubkey_map[pubkey]; !ok {
		key = genKey()
		deskey := sm2.Encrypt(pubkey, key)
		keys := [2][]byte{key, deskey}
		pubkey_map[pubkey] = keys

		//pubkey_map[pubkey]["des_cipher"] = deskey
	} else {
		key = keys[0]
	}
	sm4.EncryptSetKey(key[:])
	output := sm4.Encrypt(([]byte)(text))
	des_cipyer := pubkey_map[pubkey][1]
	ciphertext := hex.EncodeToString(output)
	//keyText := hex.EncodeToString(key)

	return &Envelope{
		Deskey:     hex.EncodeToString(des_cipyer),
		Ciphertext: ciphertext,
	}
}

func EncrpytEnvelopSafe(text, pubkey string) *Envelope {
	var sm4 = new(SM4)
	var key []byte
	//var ok bool
	m.RLock()
	if keys, ok := pubkey_map[pubkey]; !ok {
		m.RUnlock()
		m.Lock()
		if keys, ok = pubkey_map[pubkey]; !ok {
			key = genKey()
			deskey := sm2.Encrypt(pubkey, key)
			keys = [2][]byte{key, deskey}
			pubkey_map[pubkey] = keys
		}
		m.Unlock()
		//pubkey_map[pubkey]["des_cipher"] = deskey
	} else {
		m.RUnlock()
		key = keys[0]
	}
	sm4.EncryptSetKey(key[:])
	output := sm4.Encrypt(([]byte)(text))
	des_cipyer := pubkey_map[pubkey][1]
	ciphertext := hex.EncodeToString(output)
	//keyText := hex.EncodeToString(key)

	return &Envelope{
		Deskey:     hex.EncodeToString(des_cipyer),
		Ciphertext: ciphertext,
	}
}
func (envelope *Envelope) DecryptEnvelop(prikey string) string {
	return Decrypt(envelope.Ciphertext, prikey, envelope.Deskey)

}
func (envelope *Envelope) DecryptEnvelopSafe(prikey string) string {
	return DecryptSafe(envelope.Ciphertext, prikey, envelope.Deskey)

}

// 解析数字信封
func Decrypt(ciphertext, prikey, deskey string) string {
	var sm4 = new(SM4)
	var randkey []byte
	var ok bool

	if randkey, ok = key_map[deskey]; !ok {
		randkey = sm2.DecryptStringOrigin(prikey, deskey)
		key_map[deskey] = randkey
	}

	ciphertextByte, _ := hex.DecodeString(ciphertext)
	sm4.DecryptSetKey(randkey)
	textByte := sm4.Decrypt(ciphertextByte)
	index := bytes.IndexByte(textByte, 0)
	if index < 0 {
		return string(textByte)
	}
	trim_byte := textByte[0:index]

	return string(trim_byte)
}

func DecryptSafe(ciphertext, prikey, deskey string) string {
	var sm4 = new(SM4)
	var randkey []byte
	var ok bool
	m.RLock()
	if randkey, ok = key_map[deskey]; !ok {
		m.RUnlock()
		m.Lock()
		if randkey, ok = key_map[deskey]; !ok {
			randkey = sm2.DecryptStringOrigin(prikey, deskey)
			key_map[deskey] = randkey
		}
		m.Unlock()
	} else {
		m.RUnlock()
	}

	ciphertextByte, _ := hex.DecodeString(ciphertext)
	sm4.DecryptSetKey(randkey)
	textByte := sm4.Decrypt(ciphertextByte)
	index := bytes.IndexByte(textByte, 0)
	if index < 0 {
		return string(textByte)
	}
	trim_byte := textByte[0:index]

	return string(trim_byte)
}

func Sign(prikey string, msg []byte) string {
	return sm2.Sign(prikey, msg)
}

func VerifySign(pubkey string, msg []byte, signVaule string) int {

	sm2Verify := new(SM2)
	sm2Verify.Init()
	sm2Verify.VerifyCtx = &sm2SignVerifyContext{}
	sm2Verify.VerifyCtx.Init = 0

	verifyInitContext(sm2Verify.SM2Params, sm2Verify.VerifyCtx, pubkey)
	return signVerify(sm2Verify.SM2Params, sm2Verify.VerifyCtx, msg, signVaule)
}

//需方合并新加入
func Md5Hex(text string) string {

	h := md5.New()
	h.Write(([]byte)(text))
	c := h.Sum(nil)
	c_t := hex.EncodeToString(c)
	return c_t
}
func EncryptSm4(text string) []byte {
	var sm4 = new(SM4)
	sm4.EncryptSetKey(key_sm4[:])
	return sm4.Encrypt([]byte(text))
	//return hex.EncodeToString(output)
}
func DecryptSm4(text []byte) string {
	var sm4 = new(SM4)
	sm4.DecryptSetKey(key_sm4[:])
	text_byte := sm4.Decrypt(text)
	fmt.Println("test_byte :", text_byte)
	index := bytes.IndexByte(text_byte, 0)
	if index == -1 {
		return string(text_byte)
	}
	trim_byte := text_byte[0:index]
	return string(trim_byte)
	//return hex.EncodeToString(output)
}
