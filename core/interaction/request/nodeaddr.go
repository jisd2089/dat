package request

/**
    Author: luzequan
    Created: 2018-01-10 15:13:32
*/
type NodeAddress struct {
	MemberId    string // 节点会员ID
	Host        string // 节点Host
	Port        string // 节点Port
	URL         string // 节点url
	Priority    int    // 指定调度优先级，默认为0（最小优先级为0）
	Connectable bool   // 节点地址是否可连接
	RetryTimes  int    // 重试连接次数
}

func (self *NodeAddress) GetPriority() int {
	return self.Priority
}

func (self *NodeAddress) SetPriority(priority int) *NodeAddress {
	self.Priority = priority
	return self
}

func (a *NodeAddress) GetUrl() string {
	return "http://" + a.Host + ":" + a.Port + a.URL
}

type FileServerAddress struct {
	Host      string
	Port      int
	UserName  string
	Password  string
	TimeOut   int
	LocalDir  string
	RemoteDir string
}
