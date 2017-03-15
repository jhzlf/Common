package myList

import (
	"container/list"
	"sync"
	//	"time"
)

type MyList struct {
	lock sync.Mutex
	l    *list.List
	name string
}

func NewList(name string) *MyList {
	return &MyList{
		l:    list.New(),
		name: name}
}

func (l *MyList) Front() interface{} {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.l.Len() == 0 {
		return nil
	} else {
		s := l.l.Front()
		return s.Value
	}
}

func (l *MyList) PopFront() interface{} {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.l.Len() == 0 {
		return nil
	} else {
		s := l.l.Front()
		l.l.Remove(s)
		return s.Value
	}
}

func (l *MyList) PushBack(s interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	if s == nil {
		return false
	}
	l.l.PushBack(s)
	return true
}

func (l *MyList) FrontMoveToBack() interface{} {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.l.Len() == 0 {
		return nil
	} else {
		s := l.l.Front()
		l.l.MoveToBack(s)
		return s.Value
	}
}

func (l *MyList) Len() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.l.Len()
}

func (l *MyList) Clean() {
	l.lock.Lock()
	defer l.lock.Unlock()
	for {
		if l.l.Len() > 0 {
			l.l.Remove(l.l.Front())
		} else {
			return
		}
	}
}

//func (l *MyList) WatchLen(i time.Duration) {
//	go func(l *MyList) {
//		for {
//			log.Println("list:", l.name, "len:", l.l.Len())
//			time.Sleep(i)
//		}
//	}(l)
//}
