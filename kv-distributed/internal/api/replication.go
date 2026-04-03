package api

import (
	"encoding/json"
	"net/rpc"
	"sync"
	"time"
)

type ReplicationManager struct {
	slaves []*rpc.Client
	mu     sync.RWMutex
}

type ReplicationRequest struct {
	Operation string
	Data      []byte
	Timestamp time.Time
}

func (rm *ReplicationManager) AddSlave(address string) error {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return err
	}

	rm.mu.Lock()
	rm.slaves = append(rm.slaves, client)
	rm.mu.Unlock()
	return nil
}

func (rm *ReplicationManager) Replicate(operation string, data interface{}) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	replicationReq := ReplicationRequest{
		Operation: operation,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	replicationReq.Data = jsonData

	for _, slave := range rm.slaves {
		var resp Response
		go slave.Call("KVServer.Replicate", replicationReq, &resp)
	}
}
