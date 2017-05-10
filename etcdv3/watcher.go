package etcdv3

import (
	"fmt"

	"time"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

// watcher is the implementaion of grpc.naming.Watcher
type watcher struct {
	re            *resolver // re: Etcd Resolver
	client        etcd3.Client
	addMap        map[string]int64
	isInitialized bool
}

// Close do nothing
func (w *watcher) Close() {
	w.client.Close()
}

// Next to return the updates
func (w *watcher) Next() ([]*naming.Update, error) {
	// prefix is the etcd prefix/value to watch
	prefix := fmt.Sprintf("/%s/%s/", Prefix, w.re.serviceName)

	// check if is initialized
	if !w.isInitialized {
		// query addresses from etcd
		ctx, cancel := context.WithCancel(context.Background())
		resp, err := w.client.Get(ctx, prefix, etcd3.WithPrefix())
		w.isInitialized = true
		if err == nil {
			addrs := extractAddrs(resp)
			//if not empty, return the updates or watcher new dir
			if l := len(addrs); l != 0 {
				updates := make([]*naming.Update, l)
				for i := range addrs {
					updates[i] = &naming.Update{Op: naming.Add, Addr: addrs[i]}
				}
				cancel()
				return updates, nil
			}
		}
		cancel()
	}
	// generate etcd Watcher
	ctx, cancel := context.WithCancel(context.Background())
	rch := w.client.Watch(ctx, prefix, etcd3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				if _, ok := w.addMap[string(ev.Kv.Value)]; !ok {
					w.addMap[string(ev.Kv.Value)] = time.Now().Unix()
					cancel()
					return []*naming.Update{{Op: naming.Add, Addr: string(ev.Kv.Value)}}, nil
				}
			case mvccpb.DELETE:
				delete(w.addMap, string(ev.Kv.Value))
				cancel()
				return []*naming.Update{{Op: naming.Delete, Addr: string(ev.Kv.Value)}}, nil
			}
		}
	}
	cancel()
	return nil, nil
}

func extractAddrs(resp *etcd3.GetResponse) []string {
	addrs := []string{}

	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			addrs = append(addrs, string(v))
		}
	}

	return addrs
}
