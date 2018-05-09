package member

/**
    Author: luzequan
    Created: 2018-05-08 18:11:25
*/
var (
	memberInfos  *MemberInfoList
	partnersInfo *PartnerInfoList
)

type MemberInfoList struct {
	Head             Head              `xml:"head"`
	MemberDetailList *MemberDetailList `xml:"member_dtl_list"`
}

type PartnerInfoList struct {
	Head              Head               `xml:"head"`
	PartnerDetailList *PartnerDetailList `xml:"partner_dtl_list"`
}

type Head struct {
	FileName          string `xml:"fileName"`
	FileCreateTime    string `xml:"fileCreateTime"`
	FileCreateTimeStr string `xml:"fileCreateTimeStr"`
}

type MemberDetailList struct {
	MemberDetailInfo []*MemberDetailInfo `xml:"mem_dtl_info"`
}

type PartnerDetailList struct {
	PartnerDetailInfo []*MemberDetailInfo `xml:"partner_dtl_info"`
}

type MemberDetailInfo struct {
	MemberId  string `xml:"memId"`
	PubKey    string `xml:"pubKey"`
	SvrURL    string `xml:"svrURL"`
	Status    string `xml:"status"`
	TotLmt    string `xml:"totLmt"`
	DepBal    string `xml:"depBal"`
	SettFlag  string `xml:"settFlag"`
	SettPoint string `xml:"settPoint"`
	Threshold string `xml:"threshold"`
	AppId     string `xml:"app_id"`
	Xregcode  string `xml:"xregcode"`
}

func SetMemberInfoList(memberInfoList *MemberInfoList) *MemberInfoList {
	memberInfos = memberInfoList
	return memberInfos
}

func GetMemberInfoList() *MemberInfoList {
	return memberInfos
}

func SetPartnersInfo(partnersInfoList *PartnerInfoList) *PartnerInfoList {
	partnersInfo = partnersInfoList
	return partnersInfo
}

func GetPartnersInfo() *PartnerInfoList {
	return partnersInfo
}
