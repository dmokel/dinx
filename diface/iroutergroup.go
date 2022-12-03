package diface

// IRouterGroup ...
type IRouterGroup interface {
	DoMessageRouter(req IRequest)
	AddRouter(msgID uint32, router IRouter)
}
