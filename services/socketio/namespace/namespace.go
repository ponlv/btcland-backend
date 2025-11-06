package namespace

import (
	"fmt"
	"sync"

	"api/services/socketio/namespace/contract"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

type namespace struct {
	server      *socket.Server
	name        string
	connections map[socket.SocketId]*socket.Socket

	sync sync.Mutex
}

func NewNamespace(server *socket.Server, name string) contract.Namespace {
	return &namespace{
		server:      server,
		name:        name,
		connections: map[socket.SocketId]*socket.Socket{},
	}
}

func (n *namespace) namespace() socket.Namespace {
	return n.server.Of(fmt.Sprintf("/%s", n.name), nil)
}

func (n *namespace) Name() string {
	return n.name
}

func (n *namespace) Connections() map[socket.SocketId]*socket.Socket {
	return n.connections
}

func (n *namespace) AddConnection(conn *socket.Socket) {
	n.sync.Lock()
	n.connections[conn.Id()] = conn
	n.sync.Unlock()
}

func (n *namespace) DeleteConnection(id socket.SocketId) {
	n.sync.Lock()
	delete(n.connections, id)
	n.sync.Unlock()
}

func (n *namespace) OnConnect(fn func(*socket.Socket)) error {
	var err error
	err = n.namespace().On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		err = client.Emit("auth", client.Handshake().Auth)
		if err != nil {
			return
		}

		fn(client)
	})

	return err
}

func (n *namespace) Broadcast(event string, args ...interface{}) error {
	return n.namespace().Emit(event, args...)
}
