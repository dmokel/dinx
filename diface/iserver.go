package diface

// IServer define the method that must be implemented by Server structure
type IServer interface {
	Start()
	Stop()
	Serve()

	GetConnectionManager() IConnectionManager
	AddRouter(msgID uint32, router IRouter)
}
