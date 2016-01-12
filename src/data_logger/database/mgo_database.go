package database

import "gopkg.in/mgo.v2"

// Mongo Database Implemenation

// MgoDatabase holds connection information of a certain Mongo instance
type MgoDatabase struct {
	Host       string
	Db         string
	Collection string
	session    *mgo.Session
}

// WriteElement writes data to the Mongodb database
func (mongo *MgoDatabase) WriteElement(element Measurement) error {
	var err error
	mongo.session, err = mongo.connect()
	if err == nil {
		defer mongo.session.Close()
		mongo.writeRecord(element)
	}
	return err
}

func (mongo *MgoDatabase) connect() (*mgo.Session, error) {
	session, err := mgo.Dial(mongo.Host)
	if err == nil {
		session.SetMode(mgo.Monotonic, true)
	}
	return session, err
}

func (mongo *MgoDatabase) writeRecord(element Measurement) {
	c := mongo.session.DB(mongo.Db).C(mongo.Collection)
	err := c.Insert(element)
	if err != nil {
		return
	}
}
