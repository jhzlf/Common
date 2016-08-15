package websocket

import (
	"Common/go-engine.io/message"
	"Common/go-engine.io/parser"
	"Common/go-engine.io/transport"
	//"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

type Server struct {
	callback transport.Callback
	conn     *websocket.Conn
}

func NewServer(w http.ResponseWriter, r *http.Request, callback transport.Callback) (transport.Server, error) {
	//	fmt.Println(">>>>>>>>>>>>>>>> NewServer", r.Header, r.Body)
	conn, err := websocket.Upgrade(w, r, nil, 10240, 10240)
	if err != nil {
		return nil, err
	}

	ret := &Server{
		callback: callback,
		conn:     conn,
	}

	go ret.serveHTTP(w, r)

	return ret, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Server) NextWriter(msgType message.MessageType, packetType parser.PacketType) (io.WriteCloser, error) {
	wsType, newEncoder := websocket.TextMessage, parser.NewStringEncoder
	if msgType == message.MessageBinary {
		wsType, newEncoder = websocket.BinaryMessage, parser.NewBinaryEncoder
	}

	w, err := s.conn.NextWriter(wsType)
	if err != nil {
		return nil, err
	}
	ret, err := newEncoder(w, packetType)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *Server) Close() error {
	return s.conn.Close()
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	defer s.callback.OnClose(s)

	//fmt.Println(">>>>>>>>>>>>>>>serveHttp websocket Begin server.go:62 | ", r.Header, r.RequestURI)

	for {
		t, r, err := s.conn.NextReader()
		if err != nil {
			s.conn.Close()
			return
		}

		//fmt.Println(">>>>>>>>>>>>>>>serveHttp websocket server.go:71 | ", t, r)
		switch t {
		case websocket.TextMessage:
			fallthrough
		case websocket.BinaryMessage:
			decoder, err := parser.NewDecoder(r)
			if err != nil {
				return
			}
			s.callback.OnPacket(decoder)
			decoder.Close()
		}
	}
}
