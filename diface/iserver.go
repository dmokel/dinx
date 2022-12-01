package diface

// IServer define the method that must be implemented by Server structure
type IServer interface {
	Start()
	Stop()
	Serve()

	AddRouter(router IRouter)
}
