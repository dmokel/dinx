package dnet

import "github.com/dmokel/dinx/diface"

// BaseRouter ...
type BaseRouter struct{}

// PreHandle ...
func (r *BaseRouter) PreHandle(req diface.IRequest) {}

// Handle ...
func (r *BaseRouter) Handle(req diface.IRequest) {}

// PostHandle ...
func (r *BaseRouter) PostHandle(req diface.IRequest) {}
