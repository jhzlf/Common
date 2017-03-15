package BaseTerminal

import (
	"Common/logger"
	//	"Common/myList"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	//	"sync"
	"sync/atomic"
	"time"
)

const (
	STATUS_NULL = iota
	STATUS_CONNECTED
	STATUS_CLOSEING
)

type CBase struct {
	//	send_buff *myList.MyList
	//	cond      *sync.Cond
	sendNum int32
	linkID  uint64
	status  int
	f       interface{}
}

type TcpClient struct {
	CBase
	reconnect bool
	base      BaseTerminal
	conn      *net.TCPConn
	room      *Broadcast
}

func NewTcpClient(b bool) *TcpClient {
	//	locker := new(sync.Mutex)
	ID = ID + 1
	return &TcpClient{
		CBase: CBase{
			//			send_buff: myList.NewList("sendbuf"),
			//			cond:      sync.NewCond(locker),
			linkID: ID,
			status: STATUS_NULL,
		},
		reconnect: b,
	}
}

func (c *CBase) LinkID() uint64 {
	return c.linkID
}

func (c *CBase) Status() int {
	return c.status
}

func (c *TcpClient) Remote() string {
	return c.conn.RemoteAddr().String()
}

func (c *TcpClient) Local() string {
	return c.conn.LocalAddr().String()
}

func (c *TcpClient) ConnectTcp(host string, port int, base BaseTerminal) error {
	c.base = base
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Panic(err.Error())
		return err
	}
	go func() {
		for {
			logger.Infof("Connect remote %s", addr.String())
			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				logger.Warn(err.Error())
				if !c.reconnect {
					break
				}
				time.Sleep(time.Second * 3)
				continue
			}
			conn.SetKeepAlive(true)
			c.conn = conn
			//			ch := make(chan int)
			//			c.asyncSend(ch)
			//			<-ch
			c.status = STATUS_CONNECTED
			base.OnConnect(c)
			c.handleTcpClient()
			if !c.reconnect {
				break
			}
			time.Sleep(time.Second * 3)
		}
	}()
	return nil
}

func (c *TcpClient) notifyClose() {
	//		logger.Debug("notifyClose Status:", c.status, " Sid:", c.Sid())
	if c.status == STATUS_CONNECTED {
		c.conn.Close()
	}
	if c.status != STATUS_NULL {
		c.base.OnClose(c.f)
	}
	c.status = STATUS_NULL
	//	if c.status != STATUS_NULL {
	//		c.send_buff.Clean()
	//		c.status = STATUS_NULL
	//		c.conn.Close()
	//		c.base.OnClose(c.f)
	//	}
}

//func (c *TcpClient) asyncSend(ch chan int) {
//	go func() {
//		defer func() {
//			c.notifyClose()
//		}()

//		ch <- 1
//		for {
//			c.cond.L.Lock()
//			c.cond.Wait()
//			c.cond.L.Unlock()

//			for {
//				p := c.send_buff.PopFront()
//				if p == nil {
//					break
//				}
//				msg, ok := p.(*BuffEx)
//				if !ok {
//					logger.Error("convert msg error")
//					continue
//				}
//				if msg.Extra == nil {
//					return
//				} else {
//					_, err := c.conn.Write(msg.Buf)
//					if err != nil {
//						logger.Warnf("%s write buffer error.", c.conn.RemoteAddr().String(), err)
//						return
//					}
//				}
//			}
//		}
//	}()
//}

func (c *TcpClient) handleTcpClient() {
	defer func() {
		c.notifyClose()
	}()

	buf := make([]byte, 2048)
	var recvBuf []byte
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				logger.Info("connection is closed.", c.conn.RemoteAddr().String())
			} else {
				logger.Warn("Read Error: ", err.Error())
			}
			return
		}
		recvBuf = append(recvBuf, buf[:n]...)

		for {
			ok, data := c.base.CheckOnePackage(&recvBuf)
			if ok {
				c.base.OnDataIn(c.f, data)
			} else {
				break
			}
		}
	}
}

func (c *TcpClient) Send(pData string) bool {
	return c.SendBin([]byte(pData))
}

func (c *TcpClient) SendBin(pData []byte) bool {
	if len(pData) == 0 {
		return false
	}
	if c.status == STATUS_CONNECTED {
		atomic.AddInt32(&c.sendNum, 1)
		go func() {
			_, err := c.conn.Write(pData)
			atomic.AddInt32(&c.sendNum, -1)
			if err != nil {
				logger.Warnf("%s write buffer error.", c.conn.RemoteAddr().String(), err)
				c.conn.Close()
				c.status = STATUS_CLOSEING
			}
		}()
		return true
	} else {
		return false
	}
}

func (c *TcpClient) Close() {
	//	c.send_buff.PushBack(&BuffEx{nil, nil})
	//	c.status = STATUS_CLOSEING
	//	c.cond.Broadcast()
	c.reconnect = false
	c.status = STATUS_CLOSEING
	go func() {
		i := 0
		for {
			v := atomic.LoadInt32(&c.sendNum)
			if v == 0 || i == 5 {
				c.conn.Close()
				break
			} else {
				time.Sleep(time.Second)
				i++
			}
		}
	}()
}

func (c *TcpClient) Join(room string) {
	if c.room == nil {
		return
	}
	c.room.Join(room, strconv.FormatUint(c.linkID, 10), c)
}

func (c *TcpClient) Leave(room string) {
	if c.room == nil {
		return
	}
	c.room.Leave(room, strconv.FormatUint(c.linkID, 10))
}

func (c *TcpClient) BroadcastTo(room string, buf string) error {
	if c.room == nil {
		return errors.New("client no room")
	}
	return c.room.Send(strconv.FormatUint(c.linkID, 10), room, buf)
}

func (c *TcpClient) GetOtherMem(room string) []interface{} {
	return c.room.GetRoomMem(strconv.FormatUint(c.linkID, 10), room)
}

func (c *TcpClient) RoomMemCount(room string) int {
	if c.room == nil {
		return 0
	}
	return c.room.Count(room)
}

func (c *TcpClient) LeaveAll() {
	if c.room == nil {
		return
	}
	c.room.LeaveAll(strconv.FormatUint(c.linkID, 10))
}

func (c *TcpClient) Check(room string) bool {
	if c.room == nil {
		return false
	}
	return c.room.Check(room, strconv.FormatUint(c.linkID, 10))
}
