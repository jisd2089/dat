package transfer

/**
    Author: luzequan
    Created: 2018-08-01 10:13:09
*/
import (
	"drcs/core/interaction/response"
	"encoding/base64"
	"net/url"
)

type EncodeTransfer struct{}

func NewEncodeTransfer() Transfer {
	return &EncodeTransfer{}
}

// 封装encode服务
func (ft *EncodeTransfer) ExecuteMethod(req Request) Response {

	var (
		err        error
		body       string
		bodyByte   []byte
		returnCode = "000000"
		retMsg     = "encode or decode success"
	)

	switch req.GetMethod() {
	case "BASE64ENCODE":
		requestTxt := req.GetParameters()
		body = Base64Encode(requestTxt)
	case "BASE64DECODE":
		ciphertext := req.GetPostData()
		bodyByte, err = base64.StdEncoding.DecodeString(ciphertext)
	case "URLENCODE":
		urlstr := req.Param("urlstr")
		body, err = URLEncode(urlstr)
	case "URLDECODE":
		compressFile(req)
	}

	if err != nil {
		returnCode = "000007"
		retMsg = err.Error()
	}

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: returnCode,
		BodyStr:    body,
		Body:       bodyByte,
		ReturnMsg:  retMsg,
	}
}

// base64 编解码 #######################################################################
func Base64Encode(plaintext []byte) string {
	return base64.StdEncoding.EncodeToString(plaintext)
}

func Base64Decode(ciphertext string) ([]byte, error) {
	base64.StdEncoding.DecodeString(ciphertext)
	return nil, nil
}

// URL 编解码 #######################################################################
func URLEncode(urlstr string) (string, error) {
	urlVal := make(url.Values)
	urlVal.Add(urlstr, "")
	return urlVal.Encode()[:len(urlVal.Encode())-1], nil
}

func URLDecode(ciphertext []byte) ([]byte, error) {

	return nil, nil
}

func (ft *EncodeTransfer) Close() {}
