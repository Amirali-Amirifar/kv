package kvNode

type SnapshotEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (n *Node) CreateSnapshot() []SnapshotEntry {
	var snapshot []SnapshotEntry
	n.store.Iterate(func(key, value string) {
		snapshot = append(snapshot, SnapshotEntry{
			Key:   key,
			Value: value,
		})
	})
	return snapshot
}
