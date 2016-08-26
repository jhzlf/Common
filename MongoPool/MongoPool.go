package MongoPool

import (
	"Common/logger"
	"errors"
	//	"sync"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
)

type MongoPool struct {
	//	m_url string
	//	lock       sync.Mutex
	//	conn_type  int
	//	mgoSession *mgo.Session

}

var (
	//	b     bool
	m_map = make(map[interface{}]*mgo.Session)
)

type MongoIndex struct {
	Key        []string
	Unique     bool
	DropDups   bool
	Background bool // See notes.
	Sparse     bool
	Name       string
}

func (m *MongoPool) AddDb(name interface{}, url string) {
	//	var one sync.Once
	//	one.Do(once)
	if strings.Count(url, ":") == 0 {
		url += ":27017"
	}
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

	//	v := &MONGO_MAP{url, p}
	m_map[name] = p
	//	m_map["123"] = v
}

func getSession(name interface{}) *mgo.Session {
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

	return v.Clone()
}

func witchCollection(name interface{}, db, collection string, s func(*mgo.Collection) error) error {
	session := getSession(name)
	if session == nil {
		logger.Error("mongo session get error")
		return errors.New("mongo session get error")
	}
	defer session.Close()
	c := session.DB(db).C(collection)
	return s(c)
}

func (m *MongoPool) AddOne(name interface{}, db, collection string, p interface{}) error {
	// p.Id = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	return witchCollection(name, db, collection, query)
}

func (m *MongoPool) FindOne(name interface{}, db, collection string, i *map[string]interface{}, sort string, o interface{}) error {
	query := func(c *mgo.Collection) error {
		if sort == string("") {
			if i == nil {
				return c.Find(nil).One(o)
			} else {
				return c.Find(i).One(o)
			}
		} else {
			if i == nil {
				return c.Find(nil).Sort(sort).One(o)
			} else {
				return c.Find(i).Sort(sort).One(o)
			}
		}
	}
	return witchCollection(name, db, collection, query)
}

func (m *MongoPool) Del(name interface{}, db, collection string, i *map[string]interface{}) error {
	query := func(c *mgo.Collection) error {
		return c.Remove(i)
	}
	return witchCollection(name, db, collection, query)
}

func (m *MongoPool) EnsureIndex(name interface{}, db, collection string, index MongoIndex) error {
	mi := mgo.Index{
		Key:        index.Key,
		Unique:     index.Unique,
		DropDups:   index.DropDups,
		Background: index.Background,
		Sparse:     index.Sparse,
		Name:       index.Name,
	}

	query := func(c *mgo.Collection) error {
		return c.EnsureIndex(mi)
	}
	return witchCollection(name, db, collection, query)
}
