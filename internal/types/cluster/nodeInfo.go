package cluster

import "net"

type NodeInfo struct {
	ID            int
	ShardKey      int
	Status        NodeStatus
	Address       net.TCPAddr
	StoreNodeType StoreNodeType
}

func (n *NodeInfo) GetID() int {
	return n.ID
}

func (n *NodeInfo) GetAddress() (string, int) {
	return n.Address.IP.String(), n.Address.Port
}

func (n *NodeInfo) GetStatus() NodeStatus {
	return n.Status
}
