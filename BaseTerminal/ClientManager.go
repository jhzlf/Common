package BaseTerminal

import (
	//	"Common/logger"
	"sync"
)

type ClientManager struct {
	m_client_map  map[interface{}]interface{}
	m_client_lock sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		m_client_map: make(map[interface{}]interface{}),
	}
}

func (s *ClientManager) AddClient(linkID interface{}, base interface{}) interface{} {
	s.m_client_lock.Lock()
	defer s.m_client_lock.Unlock()

	v, ok := s.m_client_map[linkID]
	if ok {
		delete(s.m_client_map, linkID)
	}
	s.m_client_map[linkID] = base
	//	logger.Debug(">>>>>>>>>>>>>>>>>>>>>>>> ADD", linkID, "	", len(s.m_client_map))
	return v
}

func (s *ClientManager) DelClient(linkID interface{}) {
	s.m_client_lock.Lock()
	defer s.m_client_lock.Unlock()
	delete(s.m_client_map, linkID)
	//	logger.Debug(">>>>>>>>>>>>>>>>>>>>>>>> DEL", linkID, "	", len(s.m_client_map))
}

func (s *ClientManager) FindClient(linkID interface{}) interface{} {
	s.m_client_lock.Lock()
	defer s.m_client_lock.Unlock()
	v, ok := s.m_client_map[linkID]
	if !ok {
		return nil
	}
	return v
}

func (s *ClientManager) Clear() {
	s.m_client_lock.Lock()
	defer s.m_client_lock.Unlock()
	s.m_client_map = make(map[interface{}]interface{})
}

func (s *ClientManager) Count() int {
	s.m_client_lock.Lock()
	defer s.m_client_lock.Unlock()
	return len(s.m_client_map)
}

func (s *ClientManager) Range() ([]interface{}, []interface{}) {
	s.m_client_lock.Lock()
	defer s.m_client_lock.Unlock()
	var kin []interface{}
	var vin []interface{}
	for k, v := range s.m_client_map {
		kin = append(kin, k)
		vin = append(vin, v)
	}
	return kin, vin
}
