package smartsail

/**
    Author: luzequan
    Created: 2018-08-03 14:57:26
*/

/**
  地址可信度接口:http://service.91zhengxin.com/qdDcenter/getAddrTrust
   
  多头负债接口:http://service.91zhengxin.com/qdDcenter/getMultiLiabilities
   
  计算消费变量模型接口文档(查询电商数据):http://service.91zhengxin.com/qdDcenter/getConsume
   
  空号检测接口: http://service.91zhengxin.com/qdDcenter/getPhoneisempty
 */

const (
	SMARTSAIL_URL_TEST = "http://test.91zhengxin.com/jyzx/zxservice.do"
	SMARTSAIL_URL      = "http://service.91zhengxin.com/qdDcenter/getConsume"
	SMARTSAIL_CLIKEY   = "C99BC53C32584933AB080B86C5318CE7"

	SMARTSAIL_PRIVATE_KEY = `
-----BEGIN RSA PRIVATE KEY-----
MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBANVh9APN+kArMXbL
p2/88Mlt1KPj0NVDHutHdLPlBSil8B3z02ELQRebPBRGfzS0scO4KiAe9v5bhVoA
sEqZ1qcHWaCDozzm/dBYDg4zo+Bg+AuOrvpQc4abhmgvZH8TNFDCZrSrxqQtyAQv
RP9jOnpQI/hLS8WzCkVmrKQoYnxHAgMBAAECgYEAjdwrUAQ2ZUbSAbpvPKKaqi+c
eMDSa5XODnlY+ug9P8LiGeeqFhBXXAxWKtybYTzoGchsKSKs7nmF9EoU6ePQsc/2
sookQ36lxGMkM/GG3fR/sqxpMUuNpSZYzkwhvyO372gu36u7bab+/+V5EJmWEYOX
8WYwbL/Z7wmmrjEU87ECQQDu+pOc4CHQN1WA+LS+jGIV5m1XquQWH3Scz05swKyy
0vAbXwYqW/ehSCB4+qWpMOm1Ua6rF92MmTn1uTIUGFqjAkEA5JSrDlLfSB96NkbE
xV8Jp42HEHwLe5zRYAn6ZA3hdNg7//AfERebvfrlV7wOAobQgoLqBGgEliH9NoON
M9K2DQJBAMP/c014bYMNvuy2Ddcx38hCYm9CUyrpxYROae274GgRpKduOepH30LB
mxBd0bx/x03UnkLooeYMTYMAzte4Wa0CQESB8kqWt+jr1jsSNsNY2pHnLwXXx7FC
rNX155+5MUtNy53Hn+gFhV4JJleHO0OymCeliPHNIyLECRofj2Bq1LkCQQDeIg2h
+7LgCwsSHLCVx+k6b1KEjgErQG6doQRl2r6PQlCZ9JhydFox5XHpnqk/DvURz+Ev
kBO7M4njSJlXUnXi
-----END RSA PRIVATE KEY-----
`

	SMARTSAIL_PUBLIC_KEY = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBEQqL3Hr7ud7MrEvuZMAVzl8C
jQwjK/sTx5UXDc+pUV+uIOhKA0wEG3Or+rH1wddITcW89Ti5zv+ypz1jlOtvS8GJ
+unjxxW7f4tLcmaUKWNxbhmgXZ6I05Dssa67oWhmPV/f5/L2Wgk9NFwbKYJWF7jP
UccC4+dC9f1FTroh5QIDAQAB
-----END PUBLIC KEY-----
`
)

type RequestMsg struct {
	CliKey      string `json:"cliKey"`
	RequestData string `json:"data"`
}

type RequestData struct {
	Phone     string `json:"phone"`     // 手机号
	Name      string `json:"name"`      // 姓名
	StartTime string `json:"starttime"` // 开始时间
}
