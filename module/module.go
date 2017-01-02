package module

import (
    "ana/conf"
    "gopkg.in/mgo.v2"
    "time"
)

var (
    MongoSession *mgo.Session
)

func init() {
    mongoConf := conf.MainConf.Mongo

    info := &mgo.DialInfo{
        Addrs: mongoConf.Addrs,
        Timeout: 60 * time.Second,
        Database: mongoConf.Database,
        Username: mongoConf.Username,
        Password: mongoConf.Password,
        Mechanism: mongoConf.Mechanism,
        Source: mongoConf.Source,
    }

    var err error
    MongoSession, err = mgo.DialWithInfo(info)
    if err != nil {
        panic(err)
    }
}
