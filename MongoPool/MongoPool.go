package MongoPool

import (
	"Common/logger"
	"errors"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
)

type MONGO_MAP struct {
	m_url string
	//	lock       sync.Mutex
	//	conn_type  int
	mgoSession *mgo.Session
}

var (
	b     bool
	m_map map[string]*MONGO_MAP
)

func once() {
	if !b {
		m_map = make(map[string]*MONGO_MAP)
		b = true
	}
}

func Init(name, url string) {
	var one sync.Once
	one.Do(once)
	logger.Info("Init MongoPool add:", url)
	var p *mgo.Session
	var err error
	for {
		p, err = mgo.Dial(url)
		if err != nil {
			logger.Error("connect mongo error, ", url, " ", err)
			time.Sleep(3 * time.Second)
		} else {
			p.SetPoolLimit(20)
			break
		}
	}

	v := &MONGO_MAP{url, p}
	m_map[name] = v
	//	m_map["123"] = v
}

func getSession(name string) *mgo.Session {
	v, ok := m_map[name]
	if !ok {
		logger.Error("get session error,", name)
		return nil
	}

	//	if v.mgoSession == nil {
	//		go ConnectTo(name, v)

	//		for {
	//			time.Sleep(time.Second)
	//			if v.mgoSession != nil {
	//				break
	//			}
	//		}
	//	}

	return v.mgoSession.Clone()
}

func witchCollection(name, db, collection string, s func(*mgo.Collection) error) error {
	session := getSession(name)
	if session == nil {
		logger.Error("mongo session get error")
		return errors.New("mongo session get error")
	}
	defer session.Close()
	c := session.DB(db).C(collection)
	return s(c)
}

func AddOne(name, db, collection string, p interface{}) error {
	// p.Id = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	return witchCollection(name, db, collection, query)
}

func FindOne(name, db, collection string, i *map[string]interface{}, o interface{}) error {
	query := func(c *mgo.Collection) error {
		return c.Find(i).One(o)
	}
	return witchCollection(name, db, collection, query)
}

func Del(name, db, collection string, i *map[string]interface{}) error {
	query := func(c *mgo.Collection) error {
		return c.Remove(i)
	}
	return witchCollection(name, db, collection, query)
}
