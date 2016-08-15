package engineio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	//"strings"
	"sync"
	"time"

	"Common/go-engine.io/message"
	"Common/go-engine.io/parser"
	"Common/go-engine.io/transport"
)

type MessageType message.MessageType

const (
	MessageBinary MessageType = MessageType(message.MessageBinary)
	MessageText   MessageType = MessageType(message.MessageText)
)

// Conn is the connection object of engine.io.
type Conn interface {

	// Id returns the session id of connection.
	Id() string

	// Request returns the first http request when established connection.
	Request() *http.Request

	// Close closes the connection.
	Close() error

	// NextReader returns the next message type, reader. If no message received, it will block.
	NextReader() (MessageType, io.ReadCloser, error)

	// NextWriter returns the next message writer with given message type.
	NextWriter(messageType MessageType) (io.WriteCloser, error)

	CurrentName() string
}

type transportCreaters map[string]transport.Creater

func (c transportCreaters) Get(name string) transport.Creater {
	return c[name]
}

type serverCallback interface {
	configure() config
	transports() transportCreaters
	onClose(sid string)
}

type state int

const (
	stateUnknow state = iota
	stateNormal
	stateUpgrading
	stateClosing
	stateClosed
)

type serverConn struct {
	id              string
	request         *http.Request
	callback        serverCallback
	writerLocker    sync.Mutex
	transportLocker sync.RWMutex
	currentName     string
	current         transport.Server
	upgradingName   string
	upgrading       transport.Server
	state           state
	stateLocker     sync.RWMutex
	readerChan      chan *connReader
	pingTimeout     time.Duration
	pingInterval    time.Duration
	pingChan        chan bool
}

var InvalidError = errors.New("invalid transport")

func newServerConn(id string, w http.ResponseWriter, r *http.Request, callback serverCallback) (*serverConn, error) {
	transportName := r.URL.Query().Get("transport")
	creater := callback.transports().Get(transportName)
	if creater.Name == "" {
		return nil, InvalidError
	}
	ret := &serverConn{
		id:           id,
		request:      r,
		callback:     callback,
		state:        stateNormal,
		readerChan:   make(chan *connReader),
		pingTimeout:  callback.configure().PingTimeout,
		pingInterval: callback.configure().PingInterval,
		pingChan:     make(chan bool),
	}
	transport, err := creater.Server(w, r, ret)
	if err != nil {
		return nil, err
	}
	ret.setCurrent(transportName, transport)
	if err := ret.onOpen(); err != nil {
		return nil, err
	}

	// this is hehehe
	//	go ret.pingLoop()

	return ret, nil
}

func (c *serverConn) Id() string {
	return c.id
}

func (c *serverConn) Request() *http.Request {
	return c.request
}

func (c *serverConn) NextReader() (MessageType, io.ReadCloser, error) {
	if c.getState() == stateClosed {
		return MessageBinary, nil, io.EOF
	}
	ret := <-c.readerChan
	if ret == nil {
		return MessageBinary, nil, io.EOF
	}
	return MessageType(ret.MessageType()), ret, nil
}

func (c *serverConn) NextWriter(t MessageType) (io.WriteCloser, error) {
	switch c.getState() {
	case stateUpgrading:
		for i := 0; i < 30; i++ {
			time.Sleep(50 * time.Millisecond)
			if c.getState() != stateUpgrading {
				break
			}
		}
		// if c.getState() == stateUpgrading {
		// 	return nil, fmt.Errorf("upgrading")
		// }
	case stateNormal:
	default:
		return nil, io.EOF
	}
	c.writerLocker.Lock()
	ret, err := c.getCurrent().NextWriter(message.MessageType(t), parser.MESSAGE)
	if err != nil {
		c.writerLocker.Unlock()
		return ret, err
	}
	writer := newConnWriter(ret, &c.writerLocker)
	return writer, err
}

func (c *serverConn) Close() error {

	// debug.PrintStack()
	// log.Println(">>>>>>>>>>>>>>>>>>serverConn::Close ", c.Request().Header, " ", c.Request().Body)

	if c.getState() != stateNormal && c.getState() != stateUpgrading {
		return nil
	}
	if c.upgrading != nil {
		c.upgrading.Close()
	}
	c.writerLocker.Lock()
	if w, err := c.getCurrent().NextWriter(message.MessageText, parser.CLOSE); err == nil {
		writer := newConnWriter(w, &c.writerLocker)
		writer.Close()
	} else {
		c.writerLocker.Unlock()
	}
	if err := c.getCurrent().Close(); err != nil {
		return err
	}
	c.setState(stateClosing)
	return nil
}

func (c *serverConn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	log.Println(">>>>>>>>>>>>>>>>>>serverConn::ServeHTTP", r.Host, r.Header, r.Body, r.URL.String())
	transportName := r.URL.Query().Get("transport")
	if c.currentName != transportName {
		creater := c.callback.transports().Get(transportName)
		if creater.Name == "" {
			http.Error(w, fmt.Sprintf("invalid transport %s", transportName), http.StatusBadRequest)
			return
		}
		u, err := creater.Server(w, r, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//log.Println(">>>>>>>>>>>>>>>>>>serverConn::ServeHTTP set upgrading", r.Host, r.Header, r.Body)
		c.setUpgrading(creater.Name, u)
		return
	}
	c.current.ServeHTTP(w, r)
}

func (c *serverConn) OnPacket(r *parser.PacketDecoder) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(">>>>>>>>>>>>>panic log:", err)
			debug.PrintStack()
		}
	}()

	if s := c.getState(); s != stateNormal && s != stateUpgrading {
		//log.Println(">>>>>>>>>>>>>>OnPacket Return", c.getState())
		return
	}

	// log.Println(">>>>>>>>>>>>>OnPacket State:", c.getState())
	// log.Println(">>>>>>>>>>>>>OnPacket line 206", r.Type())

	switch r.Type() {
	case parser.OPEN:
	case parser.CLOSE:
		c.getCurrent().Close()
	case parser.PING:
		c.writerLocker.Lock()
		t := c.getCurrent()
		u := c.getUpgrade()
		newWriter := t.NextWriter
		if u != nil {
			if w, _ := t.NextWriter(message.MessageText, parser.NOOP); w != nil {
				w.Close()
			}
			newWriter = u.NextWriter
		}
		if w, _ := newWriter(message.MessageText, parser.PONG); w != nil {
			io.Copy(w, r)
			w.Close()
		}
		c.writerLocker.Unlock()
		fallthrough
	case parser.PONG:
		// modify by lee
		// this is hehehe
		//		c.stateLocker.RLock()
		//		if c.state != stateClosed {
		//			c.pingChan <- true
		//		}
		//		c.stateLocker.RUnlock()
		// end

	//c.pingChan <- true
	case parser.MESSAGE:
		// modify by lee
		c_state := c.getState()
		if c_state != stateClosed {
			closeChan := make(chan struct{})
			c.readerChan <- newConnReader(r, closeChan)
			<-closeChan
			close(closeChan)
		}
		// closeChan := make(chan struct{})
		// c.readerChan <- newConnReader(r, closeChan)
		// <-closeChan
		// close(closeChan)
		r.Close()
	case parser.UPGRADE:
		c.upgraded()

		//log.Println(">>>>>>>>>>>>>>>>>>> web socket UPGRADE", c.Request().Header)

		//map[Accept:[*/*] Cookie:[io=7YoEgdAnWyXb_hrHeAAd] User-Agent:[%E4%BA%B2%E5%8A%A0%E7%9B%B4%E6%92%AD/35 CFNetwork/711.2.23 Darwin/14.0.0] Accept-Language:[zh-cn] Accept-Encoding:[gzip, deflate] Connection:[keep-alive]]
		//user_agent := c.Request().Header.Get("User-Agent")
		//bIsIOS := strings.Contains(user_agent, "CFNetwork")
		//if no orgin, means mobile client
		// if bIsIOS {
		// 	log.Println(">>>>>>>>>>>>>>>>>>>>>>server_conn::OnPacket is IOS", c)
		// }
		// log.Println(">>>>>>>>>>>>>>>>>>>>>>server_conn::OnPacket", c)
		c.writerLocker.Lock()
		//		u := c.getUpgrade()
		currentWriter := c.getCurrent().NextWriter
		// if u != nil {
		// 	if w, _ := t.NextWriter(message.MessageText, parser.MESSAGE); w != nil {
		// 		w.Close()
		// 	}
		// 	newWriter = u.NextWriter
		// }
		w, _ := currentWriter(message.MessageText, parser.MESSAGE)
		w.Write([]byte("0"))
		if w != nil {
			io.Copy(w, r)
			w.Close()
		}

		c.writerLocker.Unlock()
		//}

	case parser.NOOP:
	}
}

func (c *serverConn) OnClose(server transport.Server) {
	if t := c.getUpgrade(); server == t {
		c.setUpgrading("", nil)
		t.Close()
		return
	}
	t := c.getCurrent()
	if server != t {
		return
	}
	t.Close()
	if t := c.getUpgrade(); t != nil {
		t.Close()
		c.setUpgrading("", nil)
	}
	//c.setState(stateClosed)
	// modify by lee
	c.stateLocker.Lock()
	if c.state != stateClosed {
		c.state = stateClosed
		close(c.readerChan)
		close(c.pingChan)
	}
	c.stateLocker.Unlock()
	// end
	c.callback.onClose(c.id)
}

func (s *serverConn) onOpen() error {
	upgrades := []string{}
	for name := range s.callback.transports() {
		if name == s.currentName {
			continue
		}
		upgrades = append(upgrades, name)
	}
	type connectionInfo struct {
		Sid          string        `json:"sid"`
		Upgrades     []string      `json:"upgrades"`
		PingInterval time.Duration `json:"pingInterval"`
		PingTimeout  time.Duration `json:"pingTimeout"`
	}
	resp := connectionInfo{
		Sid:          s.Id(),
		Upgrades:     upgrades,
		PingInterval: s.callback.configure().PingInterval / time.Millisecond,
		PingTimeout:  s.callback.configure().PingTimeout / time.Millisecond,
	}
	w, err := s.getCurrent().NextWriter(message.MessageText, parser.OPEN)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(w)

	//	log.Println(">>>>>>>>>>>>>>>>>Json Resp", resp)

	if err := encoder.Encode(resp); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func (c *serverConn) getCurrent() transport.Server {
	c.transportLocker.RLock()
	defer c.transportLocker.RUnlock()

	return c.current
}

func (c *serverConn) getUpgrade() transport.Server {
	c.transportLocker.RLock()
	defer c.transportLocker.RUnlock()

	return c.upgrading
}

func (c *serverConn) setCurrent(name string, s transport.Server) {
	c.transportLocker.Lock()
	defer c.transportLocker.Unlock()

	c.currentName = name
	c.current = s
}

func (c *serverConn) setUpgrading(name string, s transport.Server) {
	c.transportLocker.Lock()
	defer c.transportLocker.Unlock()

	c.upgradingName = name
	c.upgrading = s
	c.setState(stateUpgrading)
}

func (c *serverConn) upgraded() {
	c.transportLocker.Lock()

	current := c.current
	c.current = c.upgrading
	c.currentName = c.upgradingName
	c.upgrading = nil
	c.upgradingName = ""

	c.transportLocker.Unlock()

	current.Close()
	c.setState(stateNormal)
}

func (c *serverConn) getState() state {
	c.stateLocker.RLock()
	defer c.stateLocker.RUnlock()
	return c.state
}

func (c *serverConn) setState(state state) {
	c.stateLocker.Lock()
	defer c.stateLocker.Unlock()
	c.state = state
}

func (c *serverConn) pingLoop() {
	lastPing := time.Now()
	lastTry := lastPing
	for {
		now := time.Now()
		pingDiff := now.Sub(lastPing)
		tryDiff := now.Sub(lastTry)
		select {
		case ok := <-c.pingChan:
			if !ok {
				return
			}
			lastPing = time.Now()
			lastTry = lastPing
		case <-time.After(c.pingInterval - tryDiff):
			c.writerLocker.Lock()
			if w, _ := c.getCurrent().NextWriter(message.MessageText, parser.PING); w != nil {
				writer := newConnWriter(w, &c.writerLocker)
				writer.Close()
			} else {
				c.writerLocker.Unlock()
			}
			lastTry = time.Now()
		case <-time.After(c.pingTimeout - pingDiff):
			if c.CurrentName() == "websocket" {
				log.Println(">>>>>>>>>>>>>>>>>ping loop overtime!! pingTimeout:", c.pingTimeout, "  lastPing:", lastPing, "  pingInterval:", c.pingInterval, " lastTry:", lastTry, " now", now)
				c.Close()
				return
			}
			//adapt to polling has not response heartbeat
			lastPing = time.Now()
			lastTry = lastPing
			continue
		}
	}
}

func (c *serverConn) CurrentName() string {
	return c.currentName
}
