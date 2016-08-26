package BaseTerminal

import (
	"Common/logger"
	"fmt"
	"net"
)

func ListenTcp(port int, base BaseTerminal) error {
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	go func() {
		defer func() {
			logger.Info("Close listen address: %s", listener.Addr().String())
			//			server.Close()
			listener.Close()
		}()

		logger.Info("listen address: ", listener.Addr().String())
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				continue
			}
			//			logger.Infof("%d accept new connect, remote address: %s.", port, conn.RemoteAddr().String())
			logger.Info(port, " accept new connect, remote address: ", conn.RemoteAddr().String())

			pClient := NewTcpClient(false)
			pClient.conn = conn
			pClient.base = base
			ch := make(chan int)
			pClient.asyncSend(ch)
			<-ch
			pClient.status = STATUS_CONNECTED
			pClient.f = base.OnConnect(pClient)
			go pClient.handleTcpClient()
		}
	}()
	return nil
}
