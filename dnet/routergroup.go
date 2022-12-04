package dnet

import (
	"fmt"
	"strconv"

	"github.com/dmokel/dinx/diface"
	"github.com/dmokel/dinx/utils"
)

// RouterGroup ...
type routerGroup struct {
	Routers    map[uint32]diface.IRouter
	TaskQueue  []chan diface.IRequest
	workerSize uint32
}

var _ diface.IRouterGroup = &routerGroup{}

// NewRouterGroup ...
func newRouterGroup() diface.IRouterGroup {
	return &routerGroup{
		Routers:    make(map[uint32]diface.IRouter),
		TaskQueue:  make([]chan diface.IRequest, utils.GlobalIns.WorkerPoolSize),
		workerSize: utils.GlobalIns.WorkerPoolSize,
	}
}

// DoMessageRouter ...
func (rg *routerGroup) DoMessageRouter(req diface.IRequest) {
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
func (rg *routerGroup) AddRouter(msgID uint32, router diface.IRouter) {
	if _, ok := rg.Routers[msgID]; ok {
		panic("duplicate router" + strconv.Itoa(int(msgID)))
	}
	rg.Routers[msgID] = router
}

// SendMsgToTaskQueue ...
func (rg *routerGroup) SendMsgToTaskQueue(req diface.IRequest) {
	workerID := req.GetConnection().GetConnectionID() % rg.workerSize
	fmt.Printf("[Server] WorkerID = %d handle msgID = %d\n", workerID, req.GetMessageID())
	rg.TaskQueue[workerID] <- req
}

// StartWorkerPool ...
func (rg *routerGroup) StartWorkerPool() {
	for i := 0; i < int(rg.workerSize); i++ {
		rg.TaskQueue[i] = make(chan diface.IRequest, utils.GlobalIns.MaxWorkerTaskNum)
		go rg.startWorker(i, rg.TaskQueue[i])
	}
}

func (rg *routerGroup) startWorker(workerID int, taskQueue chan diface.IRequest) {
	fmt.Printf("[Server] Worker running, workerID = %d\n", workerID)
	for {
		select {
		case req := <-taskQueue:
			rg.DoMessageRouter(req)
		}
	}
}
