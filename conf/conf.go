package conf

import (
    "flag"
    "github.com/BurntSushi/toml"
)

type ServerConf struct {
    Addr string
}

type MongoConf struct {
    Addrs     []string
    Database  string
    Username  string
    Password  string
    Mechanism string
    Source    string
}

type mainConf struct {
    Server ServerConf
    Mongo  MongoConf
}

var (
    confFile string
    MainConf = &mainConf{}
)

func init() {
    flag.StringVar(&confFile, "conf", "/path/file.toml", "config file")
    flag.Parse()

    if confFile == "" {
        panic("need args --conf=/path/file")
    }
    toml.DecodeFile(confFile, MainConf)
}
