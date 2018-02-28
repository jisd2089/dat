package member

// MemberManager 成员管理器接口
type MemberManager interface {
	// 获取成员信息
	GetMemberInfo(memID string) *MemberInfo
}

// MemberInfo 会员信息结构体
type MemberInfo struct {
	MemID      string
	PubKey     string
	SvrURL     string
	Status     string // TODO 取值范围
	TotLmt     float64
	SettFlag   string
	SettPoint  string
	Threashold string // TODO 什么类型
}

// UpdatableMemberManager 对外提供更新方法的MemberManager接口扩展
type UpdatableMemberManager interface {
	MemberManager
	Update() error
}
