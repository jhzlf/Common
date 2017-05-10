package etcdv3

import (
	"Common/logger"
	"context"
	"fmt"
	"strconv"
	"strings"

	etcd3 "github.com/coreos/etcd/clientv3"
)

const (
	UUID_KEY        = "/seqs/snowflake-uuid"
	MACHINE_ID_MASK = 0x3FF // 10bit
)

func GetMachineID(serverName, target string) uint64 {
	var err error
	if client == nil {
		client, err = etcd3.New(etcd3.Config{
			Endpoints: strings.Split(target, ","),
		})
		if err != nil {
			logger.Warn("grpclb: creat etcd3 client failed: ", err)
			return 0
		}
	}

	var prevValue int
	for {
		resp, err := client.Get(context.Background(), UUID_KEY)
		if err != nil {
			logger.Warn("etcd get error ", err)
			return 0
		}
		for _, value := range resp.Kvs {
			prevValue, _ = strconv.Atoi(string(value.Value))
		}
		_, err = client.Put(context.Background(), UUID_KEY, fmt.Sprint(prevValue+1))
		if err != nil {
			logger.Warn("etcd put error ", err)
			continue
		}
		// record serial number of this service, already shifted
		return (uint64(prevValue+1) & MACHINE_ID_MASK) << 12
	}
}
