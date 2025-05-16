package service

import "github.com/Amirali-Amirifar/kv/internal"

type NodeInfo struct {
	internal.KvStoreConfig
}
type NodeManager struct {
	Nodes []NodeInfo
}

func NewNodeManager() *NodeManager {
	return &NodeManager{}
}

func (nm *NodeManager) GetNodeInfo(name string) (NodeInfo, error) {
	return NodeInfo{}, nil
}
