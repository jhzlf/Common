package myList

import "sync"

type OutRoutine interface {
	OnDataIn(p interface{})
}

type MsgRoutine struct {
	msgList *MyList
	out     OutRoutine
	cond    *sync.Cond
}

func (m *MsgRoutine) AddMsg(msg interface{}) {
	m.msgList.PushBack(msg)
	// cond.Broadcast()
	m.cond.Signal()
}

func CreateRoutine(o OutRoutine) *MsgRoutine {
	m := new(MsgRoutine)
	m.msgList = NewList("msg")
	m.out = o
	m.cond = sync.NewCond(new(sync.Mutex))

	go func() {
		for {
			m.cond.L.Lock()
			m.cond.Wait()

			for {
				p := m.msgList.PopFront()
				if p == nil {
					break
				}
				if m.out != nil {
					m.out.OnDataIn(p)
				}
			}
			m.cond.L.Unlock()
		}
	}()

	return m
}
