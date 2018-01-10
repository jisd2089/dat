package request

/**
    Author: luzequan
    Created: 2018-01-10 15:13:32
*/
type NodeAddress struct {
	IP       string // 节点IP
	URL      string // 节点url
	Priority int    //指定调度优先级，默认为0（最小优先级为0）
}

func (self *NodeAddress) GetPriority() int {
	return self.Priority
}

func (self *NodeAddress) SetPriority(priority int) *NodeAddress {
	self.Priority = priority
	return self
}