package BaseTerminal

import (
	"Common"
	"Common/HttpServer"
	"Common/logger"
	"encoding/base64"
	//	"Common/myList"
	"errors"
	"net/http"
	"strings"
	//	"sync"
	"sync/atomic"
	"time"

	"Common/go-socket.io"
)

const (
	MAX_SEND_BUFF      = 500
	SEND_TIME_INTERVAL = 500 //ms
)

type SocketIOListen struct {
	Port    int
	Timeout int
	Crt     string
	Key     string
}

type SocketIOClient struct {
	CBase
	base   SocketIOBase
	ss     socketio.Socket
	server *SocketIOServer
	enKey  []byte
}

type SocketIOServer struct {
	*socketio.Server
	room *broadcast
}

func (s *SocketIOServer) SocketIOBroadcastTo(room, buf string) {
	//	s.BroadcastTo(room, "msg", args...)
	// logger.Debug(">>>>>>>>>>>>>>>>>>>>>>>>>>>", room, " ", buf)
	s.room.Send("", room, buf)
}

func (s *SocketIOServer) SocketIOBroadcastToEncrypt(room, buf string) {
	s.room.SendEncrypt("", room, buf)
}

func (s *SocketIOServer) SocketIOCloseRoom(room string) {
	s.room.Close(room)
}

func newSocketIOClient() *SocketIOClient {
	//	locker := new(sync.Mutex)
	ID = ID + 1
	return &SocketIOClient{
		CBase: CBase{
			//			send_buff: myList.NewList("sendbuf"),
			//			cond:      sync.NewCond(locker),
			linkID: ID,
			status: STATUS_NULL,
		},
	}
}

func (c *SocketIOClient) notifyClose() {
	//	logger.Debug("notifyClose Status:", c.status, " Sid:", c.Sid())
	if c.status == STATUS_CONNECTED {
		c.Close()
	} else {
		//	if c.status == STATUS_CLOSEING {
		//		c.send_buff.Clean()
		c.status = STATUS_NULL
		c.ss.Close()
		c.base.OnClose(c.f)
		//		c.ss = nil
	}
	//	 else {

	//	}
	//	if c.status != STATUS_NULL {
	//		c.status = STATUS_NULL
	//		c.ss.Close()
	//		c.base.OnClose(c.f)
	//		c.ss = nil
	//	}
}

func (c *SocketIOClient) Sid() string {
	return c.ss.Id()
}

//func (c *SocketIOClient) asyncSend() {
//	go func() {
//		//		defer func() {
//		//			c.notifyClose()
//		//		}()

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
//					if c.status == STATUS_CLOSEING {
//						c.ss.Close()
//					}
//					return
//				} else {
//					i := time.Now().UnixNano()
//					err := c.ss.Emit("msg", string(msg.Buf))
//					if err != nil {
//						logger.Warn("write buffer error.", c.ss, "	", err, "	", time.Now().UnixNano()-i)
//						c.status = STATUS_CLOSEING
//						c.ss.Close()
//						return
//					}
//					//					timeNow := time.Now().UnixNano()
//					//					if timeNow-i < SEND_TIME_INTERVAL*1000*1000 {
//					//						sleepTime := i + 800*1000*1000 - timeNow
//					//						time.Sleep(time.Nanosecond * time.Duration(sleepTime))
//					//					}
//				}
//			}
//		}
//	}()
//}

func (c *SocketIOClient) Join(room string) {
	//	c.ss.Join(room)
	c.server.room.Join(room, c.Sid(), c)
}

func (c *SocketIOClient) Leave(room string) {
	//	c.ss.Leave(room)
	c.server.room.Leave(room, c.Sid())
}

func (c *SocketIOClient) Check(room string) bool {
	return c.server.room.Check(room, c.Sid())
}

func (c *SocketIOClient) BroadcastTo(room string, buf string) error {
	//	return c.ss.BroadcastTo(room, "msg", args...)
	return c.server.room.Send(c.Sid(), room, buf)
}

func (c *SocketIOClient) BroadcastToEncrypt(room string, buf string) error {
	return c.server.room.SendEncrypt(c.Sid(), room, buf)
}

func (c *SocketIOClient) RoomMemCount(room string) int {
	//	return c.ss.RoomMemCount(room)
	return c.server.room.Count(room)
}

func (c *SocketIOClient) LeaveAll() {
	c.server.room.LeaveAll(c.Sid())
}

func (c *SocketIOClient) Send(send string) bool {
	if len(send) == 0 {
		return false
	}
	if c.status == STATUS_CONNECTED {
		//		c.send_buff.PushBack(&BuffEx{len(send), []byte(send)})
		//		c.cond.Broadcast()
		atomic.AddInt32(&c.sendNum, 1)
		go func() {
			// logger.Debug(">>>>>>>>>>>>>>>>>>>>>> Emit ", c.Sid(), " ", send)
			err := c.ss.Emit("msg", send)
			atomic.AddInt32(&c.sendNum, -1)
			if err != nil {
				logger.Warn("write buffer error.", c.ss, "	", err)
				c.ss.Close()
			}
		}()
		return true
	} else {
		return false
	}
}

func (c *SocketIOClient) Close() {
	//	if c.status == STATUS_CONNECTED {
	//	c.send_buff.PushBack(&BuffEx{nil, nil})
	//	c.status = STATUS_CLOSEING
	//	c.cond.Broadcast()
	//	}
	c.status = STATUS_CLOSEING
	go func() {
		i := 0
		for {
			v := atomic.LoadInt32(&c.sendNum)
			if v == 0 || i == 5 {
				c.ss.Close()
				break
			} else {
				time.Sleep(time.Second)
				i++
			}
		}
	}()
}

//func InitSocketIOServer(port, conn_max, timeout int, base SocketIOBase, server_crt, server_key string) (*SocketIOServer, error) {
func InitSocketIOServer(listen []SocketIOListen, conn_max int, base SocketIOBase) (*SocketIOServer, error) {
	if len(listen) == 0 {
		return nil, errors.New("listen is empty")
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		logger.Panic(err)
	}
	if conn_max == 0 {
		conn_max = 300000
	}
	server.SetMaxConnection(conn_max)

	s := &SocketIOServer{server, newBroadcast()}

	server.On("connection", func(so socketio.Socket) {
		//		logger.Debug(pHttpServer.Server.Addr, " accept new connect from ", so.Conn().Request().RemoteAddr)Host
		logger.Debug(so.Conn().Request().Host, " accept new connect from ", so.Conn().Request().RemoteAddr)
		pClient := newSocketIOClient()
		pClient.ss = so
		pClient.base = base
		pClient.server = s
		//		pClient.asyncSend()
		pClient.status = STATUS_CONNECTED
		pClient.f = base.OnConnect(pClient)

		pClient.ss.On("disconnection", func() {
			//			if pClient.status == STATUS_CONNECTED {
			//				pClient.send_buff.PushBack(&BuffEx{nil, nil})
			//				pClient.status = STATUS_NULL
			//				pClient.cond.Broadcast()
			//				pClient.base.OnClose(pClient.f)
			//			} else {
			//				pClient.status = STATUS_NULL
			//				pClient.base.OnClose(pClient.f)
			//			}
			if pClient.status != STATUS_NULL {
				pClient.base.OnClose(pClient.f)
				pClient.status = STATUS_NULL
			}
		})

		pClient.ss.On("msg", func(data string) string {
			return base.OnDataIn(pClient.f, data)
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		logger.Error("error:", err)
	})

	for _, v := range listen {
		pHttpServer := &HttpServer.HttpServer{Handler: http.NewServeMux()}
		pHttpServer.Handler.Handle("/socket.io/", server)
		pHttpServer.Handler.Handle("/", http.FileServer(http.Dir("./asset")))
		if len(v.Crt) > 0 && len(v.Key) > 0 {
			pHttpServer.ListenAndServeTLS(v.Port, v.Timeout, v.Crt, v.Key)
			logger.Info("start socketIO server https ,", v.Port)
		} else {
			pHttpServer.ListenAndServe(v.Port, v.Timeout)
			logger.Info("start socketIO server http ,", v.Port)
		}
	}

	return s, nil
}

func (c *SocketIOClient) GetUserAgent() string {
	requestHttp := c.ss.Conn().Request()
	return requestHttp.Header.Get("User-Agent")
}

func (c *SocketIOClient) CurrentName() string {
	return c.ss.Conn().CurrentName()
}

func (c *SocketIOClient) GetIP() string {
	if v, ok := c.ss.Conn().Request().Header["X-Real-Ip"]; ok {
		return v[0]
	} else {
		tm := c.ss.Conn().Request().RemoteAddr
		vec := strings.Split(tm, ":")
		if len(tm) > 1 {
			return vec[0]
		}
		return ""
	}
}

func (c *SocketIOClient) SetEnKey(key []byte) {
	c.enKey = key
}

func (c *SocketIOClient) AddSaltEnKey(salt []byte) {
	if len(salt) < 24 {
		logger.Warn("salt error")
	}
	for i := 0; i < 24; i++ {
		c.enKey[i] = c.enKey[i] ^ salt[i]
	}
}

func (c *SocketIOClient) SendEncrypt(send string) bool {
	if len(send) == 0 {
		return false
	}
	if len(c.enKey) > 0 {
		return c.Send(base64.StdEncoding.EncodeToString(Common.Encrypt(c.enKey, []byte(send))))
	}
	return c.Send(send)
}

func (c *SocketIOClient) GetDecryptMsg(msg string) string {
	logger.Debug("GetDecryptMsg ", c.enKey, "	", msg)
	if len(c.enKey) > 0 {
		b, err := base64.StdEncoding.DecodeString(msg)
		if err != nil {
			logger.Warn("base decrypt error ", err)
			return msg
		}
		return string(Common.Decrypt(c.enKey, b))
	}
	return msg
}

func (c *SocketIOClient) GetEncryptMsg(msg string) string {
	logger.Debug("GetEncryptMsg ", c.enKey, "	", msg)
	if len(c.enKey) > 0 {
		return base64.StdEncoding.EncodeToString(Common.Encrypt(c.enKey, []byte(msg)))
	}
	return msg
}
