package BaseTerminal

import
//	"Common/logger"
//	"runtime/debug"

"sync"

type RoomClient interface {
	Send(send string) bool
}

type broadcast struct {
	_map  map[string]map[string]RoomClient
	_lock sync.Mutex
}

func newBroadcast() *broadcast {
	//return make(broadcast)
	return &broadcast{
		_map: make(map[string]map[string]RoomClient),
	}
}

func (b *broadcast) Join(room string, id string, socket RoomClient) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets, ok := b._map[room]
	if !ok || sockets == nil {
		sockets = make(map[string]RoomClient)
	}

	sockets[id] = socket
	b._map[room] = sockets

	return nil
}

func (b *broadcast) Leave(room string, id string) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets, ok := b._map[room]
	if !ok {
		return nil
	}
	delete(sockets, id)
	if len(sockets) == 0 {
		delete(b._map, room)
		return nil
	}
	b._map[room] = sockets
	return nil
}

func (b *broadcast) Check(room string, id string) bool {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets, ok := b._map[room]
	if !ok {
		return false
	}
	if _, ok := sockets[id]; ok {
		return true
	}
	return false
}

func (b *broadcast) Close(room string) {
	b._lock.Lock()
	defer b._lock.Unlock()
	delete(b._map, room)
}

func (b *broadcast) LeaveAll(id string) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	for k, sockets := range b._map {
		delete(sockets, id)
		if len(sockets) == 0 {
			delete(b._map, k)
			continue
		}
		b._map[k] = sockets
	}
	return nil
}

func (b *broadcast) Send(id, room, buf string) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets := b._map[room]
	// logger.Debug(">>>>>>>>>> ", len(sockets))
	//	debug.PrintStack()
	for k, s := range sockets {
		if k == id {
			// logger.Debug(">>>>>>>>>> 11111111111", k, " ", id)
			continue
		}
		// logger.Debug(">>>>>>>>>> 222222222222", k, " ", id)
		if s.Send(buf) == false {
			delete(sockets, k)
		}
	}
	b._map[room] = sockets
	return nil
}

func (b *broadcast) Count(room string) int {
	b._lock.Lock()
	defer b._lock.Unlock()

	return len(b._map[room])
}
