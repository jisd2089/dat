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
		returnCode = "000000"
	)

	switch req.GetMethod() {
	case "Base64Encode":
		requestTxt := req.GetParameters()
		body = Base64Encode(requestTxt)
	case "Base64Decode":
		ciphertext := ""
		base64.StdEncoding.DecodeString(ciphertext)
	case "URLEncode":
		urlstr := ""
		body, err = URLEncode(urlstr)
	case "URLDecode":
		compressFile(req)
	}

	if err != nil {
		returnCode = "000007"
	}

	return &response.DataResponse{
		StatusCode: 200,
		ReturnCode: returnCode,
		BodyStr:    body,
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
	uri, err := url.Parse(urlstr)
	if err != nil {
		return "", err
	}
	return uri.Query().Encode(), nil
}

func URLDecode(ciphertext []byte) ([]byte, error) {

	return nil, nil
}

func (ft *EncodeTransfer) Close() {}
