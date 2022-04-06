package diface

// IConnMgr ...
type IConnMgr interface {
	Add(conn IConn)
	Remove(conn IConn)
	Get(connID uint32) (IConn, error)
	Len() int
	ClearConn()
}
