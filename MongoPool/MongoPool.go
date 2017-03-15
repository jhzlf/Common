package MongoPool

import (
	"Common/logger"
	"errors"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoPool struct {
	//	m_url string
	//	lock       sync.Mutex
	//	conn_type  int
	//	mgoSession *mgo.Session

}

type QueryCondition struct {
	Name      string
	Condition interface{}
	Symbol    string
}

var (
	m_lock sync.Mutex
	m_map  = make(map[interface{}]*mgo.Session)
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
	m_lock.Lock()
	defer m_lock.Unlock()
	m_map[name] = p
}

func getSession(name interface{}) *mgo.Session {
	m_lock.Lock()
	defer m_lock.Unlock()
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

func (m *MongoPool) FindOne(name interface{}, db, collection string, i map[string]interface{}, sort string, o interface{}) error {
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

func (m *MongoPool) DelOne(name interface{}, db, collection string, i map[string]interface{}) error {
	query := func(c *mgo.Collection) error {
		return c.Remove(i)
	}
	return witchCollection(name, db, collection, query)
}

func (m *MongoPool) DelAll(name interface{}, db, collection string, i map[string]interface{}) error {
	query := func(c *mgo.Collection) error {
		_, err := c.RemoveAll(i)
		return err
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

func (m *MongoPool) Find(name interface{}, db, collection string, i map[string]interface{}, sort string, o interface{}, count *int) error {
	query := func(c *mgo.Collection) error {
		var q *mgo.Query
		if i == nil {
			q = c.Find(nil)
		} else {
			q = c.Find(i)
		}
		ct, err := q.Count()
		if err != nil {
			return err
		}
		*count = ct
		if sort == string("") {
			return q.All(o)
		} else {
			return q.Sort(sort).All(o)
		}
	}
	return witchCollection(name, db, collection, query)
}

func (m *MongoPool) FindLimit(name interface{}, db, collection string, i map[string]interface{}, sort string, limit int, o interface{}, count *int) error {
	query := func(c *mgo.Collection) error {
		var q *mgo.Query
		if i == nil {
			q = c.Find(nil)
		} else {
			q = c.Find(i)
		}
		ct, err := q.Count()
		if err != nil {
			return err
		}
		*count = ct
		if sort == string("") {
			return q.Limit(limit).All(o)
		} else {
			return q.Sort(sort).Limit(limit).All(o)
		}
	}
	return witchCollection(name, db, collection, query)
}

func (m *MongoPool) Update(name interface{}, db, collection string, s map[string]interface{}, u map[string]interface{}, mod string) error {
	switch mod {
	case "set":
		u = bson.M{"$set": u}
	default:
		logger.Warn("update mod error ", mod)
		return errors.New("update mod error")
	}
	query := func(c *mgo.Collection) error {
		return c.Update(s, u)
	}
	return witchCollection(name, db, collection, query)
}

func (m *MongoPool) MakeQueryConditionAnd(q ...*QueryCondition) map[string]interface{} {
	ret := make(map[string]interface{})
	for _, v := range q {
		if v == nil {
			continue
		}
		switch v.Symbol {
		case "=":
			ret[v.Name] = v.Condition
		case "!=":
			ret[v.Name] = bson.M{"$ne": v.Condition}
		case ">":
			ret[v.Name] = bson.M{"$gt": v.Condition}
		case "<":
			ret[v.Name] = bson.M{"$lt": v.Condition}
		case ">=":
			ret[v.Name] = bson.M{"$gte": v.Condition}
		case "<=":
			ret[v.Name] = bson.M{"$lte": v.Condition}
		case "in": //condition like this  []string{"aaa", "bbb"}
			ret[v.Name] = bson.M{"$in": v.Condition}
		}
	}
	return ret
}

func (m *MongoPool) MakeQueryConditionOr(q ...map[string]interface{}) map[string]interface{} {
	var tm []map[string]interface{}
	for _, v := range q {
		if v == nil {
			continue
		}
		tm = append(tm, v)
	}
	return bson.M{"$or": tm}
}
