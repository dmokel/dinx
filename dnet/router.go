package dnet

import "github.com/dmokel/dinx/diface"

// BaseRouter ...
type BaseRouter struct{}

var _ diface.IRouter = &BaseRouter{}

// PreHandle ...
func (r *BaseRouter) PreHandle(req diface.IRequest) {}

// Handle ...
func (r *BaseRouter) Handle(req diface.IRequest) {}

// PostHandle ...
func (r *BaseRouter) PostHandle(req diface.IRequest) {}
