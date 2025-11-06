package contract

import (
	"net/http"

	"api/services/socketio/namespace/contract"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

type Server interface {
	Server() *socket.Server
	AddNamespace(name string) contract.Namespace
	GetNamespace(name string) contract.Namespace
	GetNamespaceConnections(name string) map[socket.SocketId]*socket.Socket
	GetNamespaceConnection(name string, id socket.SocketId) *socket.Socket
	SetNamespaceConnection(name string, id socket.SocketId, conn *socket.Socket)
	RemoveNamespaceConnection(name string, id socket.SocketId)
	GetNamespaceConnectionsCount(name string) int
	RemoveNamespace(name string)
	Serve() http.Handler
	Close()
}
