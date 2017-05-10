package grpcCreate

import (
	"Common"
	"Common/etcdv3"
	"Common/logger"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type (
	KeepAlive struct {
		MaxIdle     time.Duration
		MaxAge      time.Duration
		MaxAgeGrace time.Duration
		Ping        time.Duration
		PingTimeout time.Duration
		PingMinTime time.Duration
	}

	ServerFun func(s *grpc.Server) interface{}

	ServerInf struct {
		Name  string
		Proto string
		Port  int
		Cert  string
		Key   string
		F     ServerFun
		Keep  KeepAlive
	}

	ConnectorInf struct {
		Name      string
		ServerAdd string
		Ca        string
		SName     string
		Keep      KeepAlive
	}
)

var (
	resolverMapLock sync.Mutex
	resolverMap     = make(map[string]*etcdv3.Master)
)

func CreateServer(inf []*ServerInf, path ...string) ([]interface{}, error) {
	inIP := Common.GetFirstInternal()
	// rpcReg := make(map[string]string)
	var ret []interface{}
	for _, v := range inf {
		lis, err := net.Listen(v.Proto, fmt.Sprintf(":%d", v.Port))
		if err != nil {
			logger.Error("failed to listen ", err)
			continue
		}
		var opts []grpc.ServerOption
		if v.Cert != string("") &&
			v.Key != string("") {
			creds, err := credentials.NewServerTLSFromFile(v.Cert, v.Key)
			if err != nil {
				logger.Error("Failed to generate credentials ", err)
				continue
			}
			opts = []grpc.ServerOption{grpc.Creds(creds)}
		}
		var serverKeep keepalive.ServerParameters
		if v.Keep.MaxIdle > 0 {
			serverKeep.MaxConnectionIdle = v.Keep.MaxIdle
		}
		if v.Keep.MaxAge > 0 {
			serverKeep.MaxConnectionAge = v.Keep.MaxAge
		}
		if v.Keep.MaxAgeGrace > 0 {
			serverKeep.MaxConnectionAgeGrace = v.Keep.MaxAgeGrace
		}
		if v.Keep.Ping > 0 {
			serverKeep.Time = v.Keep.Ping
		}
		if v.Keep.PingTimeout > 0 {
			serverKeep.Timeout = v.Keep.Ping
		}
		opts = append(opts, grpc.KeepaliveParams(serverKeep))

		if v.Keep.PingMinTime > 0 {
			opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				v.Keep.PingMinTime,
				true,
			}))
		}

		grpcServer := grpc.NewServer(opts...)
		ret = append(ret, v.F(grpcServer))
		logger.Info("start grpc ", v.Proto, " ", v.Name, " ", v.Port)
		go grpcServer.Serve(lis)

		if v.Name == "inside" && len(path) == 2 {
			rpcAdd := inIP + ":" + strconv.Itoa(v.Port)
			err = etcdv3.Register(path[0], rpcAdd, path[1], time.Second*10, 15)
			if err != nil {
				logger.Error("register to etcd error ", err)
				continue
			}
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
			go func() {
				s := <-ch
				logger.Info("receive signal ", s)
				etcdv3.UnRegister()
				os.Exit(1)
			}()

			// rpcReg[v.Name] = inIP + ":" + strconv.Itoa(v.Port)
		}
	}

	// if len(rpcReg) > 0 {
	// 	err := etcdv3.NewWorker(name, path, rpcReg, etcd)
	// 	if err != nil {
	// 		logger.Warn("etcd NewWorker error ", err)
	// 		return err
	// 	}
	// 	ch := make(chan os.Signal, 1)
	// 	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	// 	go func() {
	// 		s := <-ch
	// 		logger.Info("receive signal ", s)
	// 		etcdv3.DelWorker(name, path)
	// 		os.Exit(1)
	// 	}()
	// }
	return ret, nil
}

func CreateConnector(inf *ConnectorInf) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if inf.SName != string("") {
		var creds credentials.TransportCredentials
		if inf.Ca != string("") {
			var err error
			creds, err = credentials.NewClientTLSFromFile(inf.Ca, inf.SName)
			if err != nil {
				logger.Error("Failed to create TLS credentials ", err)
				return nil, err
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, inf.SName)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	if inf.Keep.Ping > 0 &&
		inf.Keep.PingTimeout > 0 {
		opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			inf.Keep.Ping,
			inf.Keep.PingTimeout,
			true,
		}))
	}

	opts = append(opts, grpc.WithTimeout(5*time.Second))
	conn, err := grpc.Dial(inf.ServerAdd, opts...)
	if err != nil {
		logger.Error("fail to dial: ", err)
		return nil, err
	}
	return conn, nil
}

func GetResolver(name, etcd string) *grpc.ClientConn {
	// resolverMapLock.Lock()
	// defer resolverMapLock.Unlock()
	// var grpcIP string
	// if master, ok := resolverMap[name]; ok {
	// 	AccessServerList := master.Members()
	// 	for _, v := range AccessServerList {
	// 		for _, a := range v.IP {
	// 			grpcIP = a
	// 		}
	// 	}
	// }
	logger.Debug("GetResolver ", name, " ", etcd)
	r := etcdv3.NewResolver(name)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, etcd, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		return nil
	}
	return conn
}
