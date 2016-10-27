package websocket

import (
	"io"
	"net/http"

	"Common/go-engine.io/message"
	"Common/go-engine.io/parser"
	"Common/go-engine.io/transport"
	"github.com/gorilla/websocket"
	// "golang.org/x/net/websocket"
)

type client struct {
	conn *websocket.Conn
	resp *http.Response
}

func NewClient(r *http.Request) (transport.Client, error) {
	dialer := websocket.DefaultDialer

	conn, resp, err := dialer.Dial(r.URL.String(), r.Header)
	if err != nil {
		return nil, err
	}

	return &client{
		conn: conn,
		resp: resp,
	}, nil
}

func (c *client) Response() *http.Response {
	return c.resp
}

func (c *client) NextReader() (*parser.PacketDecoder, error) {
	var reader io.Reader
	for {
		t, r, err := c.conn.NextReader()
		if err != nil {
			return nil, err
		}
		switch t {
		case websocket.TextMessage:
			fallthrough
		case websocket.BinaryMessage:
			reader = r
			return parser.NewDecoder(reader)
		}
	}
}

func (c *client) NextWriter(msgType message.MessageType, packetType parser.PacketType) (io.WriteCloser, error) {
	wsType, newEncoder := websocket.TextMessage, parser.NewStringEncoder
	if msgType == message.MessageBinary {
		wsType, newEncoder = websocket.BinaryMessage, parser.NewBinaryEncoder
	}

	w, err := c.conn.NextWriter(wsType)
	if err != nil {
		return nil, err
	}
	ret, err := newEncoder(w, packetType)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
