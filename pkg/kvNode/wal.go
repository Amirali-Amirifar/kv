package kvNode

import (
	"sync"
)

type WALRecord struct {
	Operation string
	Key       string
	Value     string
	Seq       int64
}

type WAL struct {
	ShardKey  int
	Records   []WALRecord
	mu        sync.RWMutex
	seq       int64
	followers map[int]int64 // Map of follower ID to their last applied sequence
}

func NewWAL(shardKey int) *WAL {
	return &WAL{
		ShardKey:  shardKey,
		Records:   make([]WALRecord, 0),
		seq:       0,
		followers: make(map[int]int64),
	}
}

func (w *WAL) Append(op, key, value string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.seq++
	record := WALRecord{
		Operation: op,
		Key:       key,
		Value:     value,
		Seq:       w.seq,
	}
	w.Records = append(w.Records, record)
}

func (w *WAL) GetLastSeq() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.seq
}

func (w *WAL) GetSince(seq int64) []WALRecord {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if seq >= w.seq {
		return nil
	}

	// Find the first record with sequence number greater than seq
	start := 0
	for i, record := range w.Records {
		if record.Seq > seq {
			start = i
			break
		}
	}

	return w.Records[start:]
}

func (w *WAL) ClearUntil(seq int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	var idx int
	for i, r := range w.Records {
		if r.Seq > seq {
			break
		}
		idx = i + 1
	}
	w.Records = w.Records[idx:]
}

func (w *WAL) UpdateFollowerProgress(followerID int, seq int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.followers[followerID] = seq
}

func (w *WAL) GetMinFollowerSeq() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.followers) == 0 {
		return 0
	}

	minSeq := w.seq
	for _, seq := range w.followers {
		if seq < minSeq {
			minSeq = seq
		}
	}
	return minSeq
}

func (w *WAL) RemoveFollower(followerID int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.followers, followerID)
}
