package diface

// IServer ...
type IServer interface {
	Start()
	Stop()
	Serve()
}
