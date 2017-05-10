package etcdv3

import (
	"Common/logger"
	"context"
	"encoding/json"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type RPCAdd map[string]string

// workerInfo is the service register information to etcd
type WorkerInfo struct {
	Name   string
	IP     RPCAdd
	CPU    int
	cancel context.CancelFunc
}

var (
	keyManager     = make(map[string]*WorkerInfo)
	keyManagerLock sync.Mutex
)

func NewWorker(name, path string, ip RPCAdd, target string) error {
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
			return err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	wInf := &WorkerInfo{
		name,
		ip,
		runtime.NumCPU(),
		cancel,
	}
	keyManagerLock.Lock()
	defer keyManagerLock.Unlock()
	key := name + "_" + path
	keyManager[key] = wInf
	go heartBeat(key, wInf, ctx)
	return nil
}

func DelWorker(name, path string) {
	key := name + "_" + path
	client.Delete(context.TODO(), key)
	keyManagerLock.Lock()
	defer keyManagerLock.Unlock()
	if v, ok := keyManager[key]; ok {
		v.cancel()
	}
}

func heartBeat(key string, info *WorkerInfo, ctx context.Context) {
	value, err := json.Marshal(info)
	if err != nil {
		logger.Error(err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			resp, err := client.Grant(context.TODO(), 10)
			if err != nil {
				logger.Error(err)
				return
			}
			_, err = client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
			if err != nil {
				logger.Error("Error: cannot put to etcd:", err)
				return
			}
		}
	}
}
