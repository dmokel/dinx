package diface

// IRouter ...
type IRouter interface {
	PreHandle(req IRequest)
	Handle(req IRequest)
	PostHandle(req IRequest)
}
