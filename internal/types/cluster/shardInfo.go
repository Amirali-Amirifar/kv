package cluster

type ShardInfo struct {
	ShardKey  int
	Master    *NodeInfo
	Followers []*NodeInfo
}

func (s *ShardInfo) GetMaster() *NodeInfo {
	return s.Master
}

func (s *ShardInfo) GetFollowers() []*NodeInfo {
	followers := make([]*NodeInfo, len(s.Followers))
	for i, f := range s.Followers {
		followers[i] = f
	}
	return followers
}
