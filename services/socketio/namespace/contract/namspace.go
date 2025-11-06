package contract

import "github.com/zishang520/socket.io/servers/socket/v3"

type Namespace interface {
	Name() string
	Connections() map[socket.SocketId]*socket.Socket

	AddConnection(conn *socket.Socket)
	DeleteConnection(id socket.SocketId)
	OnConnect(fn func(*socket.Socket)) error
	Broadcast(event string, args ...interface{}) error
}
