package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"dds/cncrypt"
	"dds/errors"
	yagrusLog "dds/log"
	"dds/settings"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

const (
	settings_xpath_filepath = "KeysFile"
)

var (
	_privateKey string
)

func GetPrivateKey() (string, error) {
	if _privateKey == "" {
		return "", fmt.Errorf("未初始化私钥")
	}
	return _privateKey, nil
}

func Initialize() error {
	//settings := settings.GetSettings()
	//filePath, err := settings.GetString(settings_xpath_filepath)
	settings := settings.GetCommomSettings()
	filePath := settings.KeysFile
	//if err != nil {
	//return fmt.Errorf("get %s from setting err:%s", settings_xpath_filepath, err.Error())
	//}

	if filePath == "" {
		return fmt.Errorf("配置缺失:%s", settings_xpath_filepath)
	}

	key, err := parseConfigFileAndCalcPriKey(filePath)
	if err != nil {
		return err
	}
	yagrusLog.Info("security init pkey:%s", key)
	_privateKey = key
	fmt.Printf("_privateKey is %s \n", _privateKey)
	// 调用国密模块的初始化方法
	cncrypt.Init(_privateKey)
	return nil
}

// 解析配置文件，并且从中计算得到私钥
func parseConfigFileAndCalcPriKey(filePath string) (string, error) {
	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("security readding file:%s error:%s", filePath, err.Error())
	}

	var memberKeys memberKeys
	err = xml.Unmarshal(text, &memberKeys)
	if err != nil {
		return "", fmt.Errorf("security unmarshall %s error:%s", filePath, err.Error())
	}
	return calculatePrivateKey(&memberKeys)
}

// 从bobo那搬来的

func calculatePrivateKey(memberKeys *memberKeys) (string, error) {
	memId := memberKeys.MemId
	userkey := memberKeys.Userkey
	prikey := memberKeys.Prikey

	hash := sha256.New()
	hash.Write([]byte(memId))
	md := hash.Sum(nil)
	key := hex.EncodeToString(md)
	key = key[len(key)-16:]

	key1 := []byte(key)
	ciphertext, _ := hex.DecodeString(userkey)
	block, _ := aes.NewCipher(key1)
	mode := cipher.NewCBCDecrypter(block, key1)
	mode.CryptBlocks(ciphertext, ciphertext)

	key2 := ciphertext
	ciphertext2, _ := hex.DecodeString(prikey)
	block, _ = aes.NewCipher(key2)
	mode = cipher.NewCBCDecrypter(block, key2)
	mode.CryptBlocks(ciphertext2, ciphertext2)
	priKeyStr := string(ciphertext2)
	if len(prikey) == 64 {
		priKeyStr = hex.EncodeToString(ciphertext2)
	}
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]+$", priKeyStr); !ok {
		return "", errors.RawNew("021000", "values in memkey.xml are not correct!")
	}
	return priKeyStr, nil
}

type memberKeys struct {
	MemId   string `xml:"memId"`
	Keyseed string `xml:"keyseed"`
	Pubkey  string `xml:"pubkey"`
	Prikey  string `xml:"prikey"`
	Userkey string `xml:"userkey"`
}

func SaveDataToxml(seed, memId, userkey string) (string, *errors.MeanfulError) {
	if memId == "" {
		return "", errors.RawNew("021001", "missing memId")
	}

	//generate private key pair
	sm2 := new(cncrypt.SM2)
	sm2.Init()
	privkey := sm2.GetPrivateKey(seed, 0)
	pubkey := sm2.GetPublicKey(privkey)
	fmt.Println(privkey, pubkey)
	//cipher private key
	hash := sha256.New()
	hash.Write([]byte(memId))
	md := hash.Sum(nil)
	digest := hex.EncodeToString(md)
	keystring := digest[len(digest)-16:]

	if userkey == "" {
		userkey = hex.EncodeToString(md[0:16])
	}

	key := []byte(keystring)
	ciphertext, _ := hex.DecodeString(userkey)
	block, _ := aes.NewCipher(key)
	mode := cipher.NewCBCDecrypter(block, key)
	mode.CryptBlocks(ciphertext, ciphertext)

	key1 := ciphertext
	ciphertext1 := []byte(privkey)
	block1, _ := aes.NewCipher(key1)
	mode = cipher.NewCBCEncrypter(block1, key1)
	mode.CryptBlocks(ciphertext1, ciphertext1)
	privkey = hex.EncodeToString(ciphertext1)

	//storage info into xml format
	header := []byte(xml.Header)
	rootbegin := "<member_keys>" + "\n"
	xmldata := append(header, rootbegin...)
	fileNamexml := "\t" + "<fileName>memkeys.xml</fileName>" + "\n"
	xmldata = append(xmldata, fileNamexml...)
	createTime := time.Now().Format("2006-01-02 15:04:05.999999")
	fileCreateTimexml := "\t" + "<fileCreateTime>" + createTime + "</fileCreateTime>" + "\n"
	xmldata = append(xmldata, fileCreateTimexml...)
	memIdxml := "\t" + "<memId>" + memId + "</memId>" + "\n"
	xmldata = append(xmldata, memIdxml...)
	pubkeyxml := "\t" + "<pubkey>" + pubkey + "</pubkey>" + "\n"
	xmldata = append(xmldata, pubkeyxml...)
	privkeyxml := "\t" + "<prikey>" + privkey + "</prikey>" + "\n"
	xmldata = append(xmldata, privkeyxml...)
	userkeyxml := "\t" + "<userkey>" + userkey + "</userkey>" + "\n"
	xmldata = append(xmldata, userkeyxml...)
	rootend := "</member_keys>"
	xmldata = append(xmldata, rootend...)

	//xmlPath, _ := settings.GetSettings().GetString("KeysFile")
	xmlPath := settings.GetCommomSettings().KeysFile
	ioutil.WriteFile(xmlPath, xmldata, os.ModeAppend)

	return pubkey, nil
}
