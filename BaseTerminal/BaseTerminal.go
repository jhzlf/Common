package BaseTerminal

var ID uint64

type BuffEx struct {
	Extra interface{}
	Buf   []byte
}

type BaseTerminal interface {
	OnConnect(pClient *TcpClient) interface{}
	OnClose(p interface{})
	OnDataIn(p interface{}, pRecv *BuffEx)
	CheckOnePackage(p *[]byte) (bool, *BuffEx)
}

type SocketIOBase interface {
	OnConnect(pClient *SocketIOClient) interface{}
	OnClose(p interface{})
	OnDataIn(p interface{}, recv string) string
}
