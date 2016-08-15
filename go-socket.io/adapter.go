package socketio

import (
	"sync"
)

//	"log"
//	"runtime/debug"

// BroadcastAdaptor is the adaptor to handle broadcasts.
type BroadcastAdaptor interface {

	// Join causes the socket to join a room.
	Join(room string, socket Socket) error

	// Leave causes the socket to leave a room.
	Leave(room string, socket Socket) error

	// Send will send an event with args to the room. If "ignore" is not nil, the event will be excluded from being sent to "ignore".
	Send(ignore Socket, room, event string, args ...interface{}) error

	// Get Room Mems
	MemCount(room string) int
}

var newBroadcast = newBroadcastDefault

//type broadcast map[string]map[string]Socket

type broadcast struct {
	_map  map[string]map[string]Socket
	_lock sync.Mutex
}

func newBroadcastDefault() BroadcastAdaptor {
	//return make(broadcast)
	return &broadcast{
		_map: make(map[string]map[string]Socket),
	}
}

func (b *broadcast) Join(room string, socket Socket) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets, ok := b._map[room]
	if !ok {
		sockets = make(map[string]Socket)
	}
	sockets[socket.Id()] = socket
	b._map[room] = sockets
	return nil
}

func (b *broadcast) Leave(room string, socket Socket) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets, ok := b._map[room]
	if !ok {
		return nil
	}
	delete(sockets, socket.Id())
	if len(sockets) == 0 {
		delete(b._map, room)
		return nil
	}
	b._map[room] = sockets
	return nil
}

func (b *broadcast) Send(ignore Socket, room, event string, args ...interface{}) error {
	//	debug.PrintStack()
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets := b._map[room]
	//log.Println(">>>>>>>>>>>>>>>>room:", room, "  count:", len(b._map[room]), &b)
	//	log.Println(">>>>>>>>>>>>>>>>room:", room, "  count:", len(b._map[room]))
	for id, s := range sockets {
		if ignore != nil && ignore.Id() == id {
			continue
		}
		s.Emit(event, args...)
	}
	return nil
}

func (b *broadcast) MemCount(room string) int {
	b._lock.Lock()
	defer b._lock.Unlock()

	//	debug.PrintStack()
	return len(b._map[room])
}
