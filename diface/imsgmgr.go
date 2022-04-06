package diface

// IMsgHandler ...
type IMsgHandler interface {
	DoMsgHandler(req IRequest)
	AddRouter(msgID uint32, router IRouter)
	StartWorkPool()
	SendMsgToTaskQueue(req IRequest)
}
