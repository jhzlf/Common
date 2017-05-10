package etcdv3

import (
	"Common/logger"
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// Master is a server
type Master struct {
	membersLock sync.Mutex
	members     map[string]*WorkerInfo
	watchPrefix string
}

func NewMaster(watchPrefix, target string) (master *Master, err error) {
	if client == nil {
		endpoints := strings.Split(target, ",")
		cfg := clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: time.Second,
		}
		var err error
		client, err = clientv3.New(cfg)
		if err != nil {
			logger.Error("Error: cannot connec to etcd:", err)
			return nil, err
		}
	}
	master = &Master{
		members:     make(map[string]*WorkerInfo),
		watchPrefix: watchPrefix,
	}
	go master.watchWorkers()
	return
}

func (m *Master) Members() (ms map[string]*WorkerInfo) {
	m.membersLock.Lock()
	defer m.membersLock.Unlock()
	ms = m.members
	return
}

func (m *Master) addWorker(key string, info *WorkerInfo) {
	m.membersLock.Lock()
	defer m.membersLock.Unlock()
	m.members[key] = info
}

func (m *Master) deleteWorker(key string) {
	m.membersLock.Lock()
	defer m.membersLock.Unlock()
	delete(m.members, key)
}

func (m *Master) watchWorkers() {
	rch := client.Watch(context.Background(), m.watchPrefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type.String() == "EXPIRE" {
				m.deleteWorker(string(ev.Kv.Key))
			} else if ev.Type.String() == "PUT" {
				info := &WorkerInfo{}
				err := json.Unmarshal(ev.Kv.Value, info)
				if err != nil {
					logger.Error(err)
				}
				m.addWorker(string(ev.Kv.Key), info)
			} else if ev.Type.String() == "DELETE" {
				m.deleteWorker(string(ev.Kv.Key))
			}
		}
	}
}
