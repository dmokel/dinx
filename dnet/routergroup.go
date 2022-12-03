package dnet

import (
	"fmt"
	"strconv"

	"github.com/dmokel/dinx/diface"
)

// RouterGroup ...
type RouterGroup struct {
	Routers map[uint32]diface.IRouter
}

// NewRouterGroup ...
func NewRouterGroup() diface.IRouterGroup {
	return &RouterGroup{
		Routers: make(map[uint32]diface.IRouter),
	}
}

// DoMessageRouter ...
func (rg *RouterGroup) DoMessageRouter(req diface.IRequest) {
	router, ok := rg.Routers[req.GetMessageID()]
	if !ok {
		fmt.Printf("not match any router, msgID = %d\n", req.GetMessageID())
		return
	}

	router.PreHandle(req)
	router.Handle(req)
	router.PostHandle(req)
}

// AddRouter ...
func (rg *RouterGroup) AddRouter(msgID uint32, router diface.IRouter) {
	if _, ok := rg.Routers[msgID]; ok {
		panic("duplicate router" + strconv.Itoa(int(msgID)))
	}
	rg.Routers[msgID] = router
}
