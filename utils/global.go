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
	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskNum uint32
}

// GlobalIns ...
var GlobalIns *global

func (g *global) reload() {
	data, err := ioutil.ReadFile("conf/dinx.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(data, GlobalIns)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalIns = &global{
		Name:             "DinxDefault",
		Version:          "V0.1",
		Network:          "tcp",
		IP:               "0.0.0.0",
		Port:             8999,
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskNum: 1024,
	}

	GlobalIns.reload()
}
