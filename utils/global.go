package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dmokel/dinx/diface"
)

type global struct {
	// Server
	Server  diface.IServer
	Network string
	IP      string
	Port    int
	Name    string

	// Dinx
	Version        string
	MaxConn        int
	MaxPackageSize uint32
}

// GlobalIns ...
var GlobalIns *global

func (g *global) reload() {
	data, err := ioutil.ReadFile("conf/dinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, GlobalIns)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalIns = &global{
		Name:           "DinxDefault",
		Version:        "V0.1",
		Network:        "tcp",
		IP:             "0.0.0.0",
		Port:           8999,
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	GlobalIns.reload()
}
