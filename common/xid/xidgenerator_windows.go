package xid

/**
    Author: luzequan
    Created: 2018-08-03 09:49:52
*/
type XidGenerator struct {
	SrcAppId    string // 源节点appid
	IdType      string // idtype
	IdNo        string // 源id值
	SrcRegCode  string // 源节点regcode
	DesAppId    string // 目的节点appid
	DesXregCode string // 目的节点xregcode
	AppXidCode  string // 源节点xidcode

	XidDealer string // xid生成远程服务标记 “0”：公安三所
	XidIp     string // 公安三所ip地址
	AppKey    string //
}

func (x *XidGenerator) GenXID() (string, error) {
	return "", nil
}

func (x *XidGenerator) ConvertXID() (string, error) {
	return "", nil
}