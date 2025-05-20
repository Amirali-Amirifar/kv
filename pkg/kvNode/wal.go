package kvNode

import (
	"sync"
	"time"
)

type WAL struct {
	ShardKey int // id of partition
	mu       sync.Mutex
	records  []WALRecord
	nextSeq  int64
}

type WALRecord struct {
	Op    string // "SET" or "DELETE"
	Key   string
	Value string // optional for DELETE
	Seq   int64  // auto-incremented
	Time  time.Time
}

func NewWAL(shardKey int) *WAL {
	return &WAL{
		ShardKey: shardKey,
		records:  make([]WALRecord, 0),
		nextSeq:  1, // or 0, depending on your design
	}
}

func (w *WAL) Append(op, key, value string) WALRecord {
	w.mu.Lock()
	defer w.mu.Unlock()
	rec := WALRecord{
		Op:    op,
		Key:   key,
		Value: value,
		Seq:   w.nextSeq,
		Time:  time.Now(),
	}
	w.records = append(w.records, rec)
	w.nextSeq++
	return rec
}

func (w *WAL) GetSince(seq int64) []WALRecord {
	w.mu.Lock()
	defer w.mu.Unlock()
	var res []WALRecord
	for _, r := range w.records {
		if r.Seq > seq {
			res = append(res, r)
		}
	}
	return res
}

func (w *WAL) ClearUntil(seq int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	var idx int
	for i, r := range w.records {
		if r.Seq > seq {
			break
		}
		idx = i + 1
	}
	w.records = w.records[idx:]
}
