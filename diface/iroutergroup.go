package diface

// IRouterGroup ...
type IRouterGroup interface {
	SendMsgToTaskQueue(req IRequest)
	DoMessageRouter(req IRequest)
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()
}
