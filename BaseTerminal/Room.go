package BaseTerminal

import (
	"time"
)

import
//	"Common/logger"
//	"runtime/debug"

"sync"

type RoomClient interface {
	Send(send string) bool
}

type RoomClientEncrypt interface {
	RoomClient
	SendEncrypt(send string) bool
}

type Broadcast struct {
	_map     map[string]map[string]RoomClient
	_mapUser map[string]map[string]int64
	_lock    sync.Mutex
}

func NewBroadcast() *Broadcast {
	//return make(Broadcast)
	return &Broadcast{
		_map:     make(map[string]map[string]RoomClient),
		_mapUser: make(map[string]map[string]int64),
	}
}

func (b *Broadcast) Join(room string, id string, socket RoomClient) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets, ok := b._map[room]
	if !ok || sockets == nil {
		sockets = make(map[string]RoomClient)
	}
	sockets[id] = socket
	b._map[room] = sockets

	rooms, ok := b._mapUser[id]
	if !ok || rooms == nil {
		rooms = make(map[string]int64)
	}
	rooms[room] = time.Now().Unix()
	b._mapUser[id] = rooms

	return nil
}

func (b *Broadcast) Leave(room string, id string) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets, _ := b._map[room]
	if sockets != nil {
		delete(sockets, id)
		if len(sockets) == 0 {
			delete(b._map, room)
		} else {
			b._map[room] = sockets
		}
	}

	rooms, _ := b._mapUser[id]
	if rooms != nil {
		delete(rooms, room)
		if len(rooms) == 0 {
			delete(b._mapUser, id)
		} else {
			b._mapUser[id] = rooms
		}
	}

	return nil
}

func (b *Broadcast) Check(room string, id string) bool {
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

func (b *Broadcast) Close(room string) {
	b._lock.Lock()
	defer b._lock.Unlock()
	if sockets, ok := b._map[room]; ok {
		for k, _ := range sockets {
			if rooms, ok := b._mapUser[k]; ok {
				delete(rooms, room)
				b._mapUser[k] = rooms
			}
		}
	}

	delete(b._map, room)
}

func (b *Broadcast) LeaveAll(id string) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	// for k, sockets := range b._map {
	// 	delete(sockets, id)
	// 	if len(sockets) == 0 {
	// 		delete(b._map, k)
	// 		continue
	// 	}
	// 	b._map[k] = sockets
	// }
	if rooms, ok := b._mapUser[id]; ok {
		for r := range rooms {
			if sockets, ok := b._map[r]; ok {
				delete(sockets, id)
				if len(sockets) == 0 {
					delete(b._map, r)
				} else {
					b._map[r] = sockets
				}
			}
		}
		delete(b._mapUser, id)
	}

	return nil
}

func (b *Broadcast) Send(id, room, buf string) error {
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

func (b *Broadcast) Count(room string) int {
	b._lock.Lock()
	defer b._lock.Unlock()

	return len(b._map[room])
}

//donot change mem by this
func (b *Broadcast) GetRoomMem(id, room string) []interface{} {
	var ret []interface{}
	b._lock.Lock()
	defer b._lock.Unlock()
	sockets := b._map[room]
	for k, s := range sockets {
		if k == id {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func (b *Broadcast) GetRooms(id string) map[string]int64 {
	b._lock.Lock()
	defer b._lock.Unlock()
	if rooms, ok := b._mapUser[id]; ok {
		return rooms
	}
	return nil
}

func (b *Broadcast) SendEncrypt(id, room, buf string) error {
	b._lock.Lock()
	defer b._lock.Unlock()

	sockets := b._map[room]
	for k, s := range sockets {
		if k == id {
			continue
		}
		if se, ok := s.(RoomClientEncrypt); ok {
			if se.SendEncrypt(buf) == false {
				delete(sockets, k)
			}
		} else {
			if s.Send(buf) == false {
				delete(sockets, k)
			}
		}
	}
	b._map[room] = sockets
	return nil
}
