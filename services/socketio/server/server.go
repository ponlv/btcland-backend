package server

import (
	"net/http"
	"os"
	"time"

	"api/services/socketio/namespace"
	"api/services/socketio/namespace/contract"
	serverContract "api/services/socketio/server/contract"
	"github.com/zishang520/socket.io/servers/engine/v3/transports"
	"github.com/zishang520/socket.io/servers/socket/v3"
	"github.com/zishang520/socket.io/v3/pkg/log"
	"github.com/zishang520/socket.io/v3/pkg/types"
)

var socketIO serverContract.Server

type server struct {
	server     *socket.Server
	namespaces []contract.Namespace
	opts       *socket.ServerOptions
}

func NewSocketIOServer(opts *socket.ServerOptions) (serverContract.Server, error) {
	if socketIO != nil {
		return socketIO, nil
	}

	if opts == nil {
		opts = socket.DefaultServerOptions()
		opts.SetServeClient(true)
		opts.SetConnectionStateRecovery(&socket.ConnectionStateRecovery{})
		opts.SetAllowEIO3(true)
		opts.SetPingInterval(25 * time.Second)
		opts.SetPingTimeout(60 * time.Second)
		opts.SetMaxHttpBufferSize(1000000)
		opts.SetConnectTimeout(1000 * time.Millisecond)
		opts.SetTransports(types.NewSet[transports.TransportCtor](socket.WebSocket, socket.Polling))
		opts.SetCors(&types.Cors{
			Origin:      "*",
			Credentials: true,
		})
	}

	log.DEBUG = os.Getenv("ENV") == "DEV"
	socketio := socket.NewServer(nil, nil)

	socketIO = &server{
		server:     socketio,
		namespaces: make([]contract.Namespace, 0),
		opts:       opts,
	}

	if err := namespace.DigitalSignNamespace(socketIO); err != nil {
		return nil, err
	}

	return socketIO, nil
}

func Server() serverContract.Server {
	return socketIO
}

func (s *server) Server() *socket.Server {
	return s.server
}

func (s *server) AddNamespace(name string) contract.Namespace {
	np := s.GetNamespace(name)
	if np != nil {
		return np
	}

	np = namespace.NewNamespace(s.server, name)
	s.namespaces = append(s.namespaces, np)

	return np
}

func (s *server) GetNamespace(name string) contract.Namespace {
	for _, ns := range s.namespaces {
		if ns.Name() == name {
			return ns
		}
	}

	return nil
}

func (s *server) GetNamespaceConnections(name string) map[socket.SocketId]*socket.Socket {
	ns := s.GetNamespace(name)
	if ns == nil {
		return nil
	}
	return ns.Connections()
}

func (s *server) GetNamespaceConnection(name string, id socket.SocketId) *socket.Socket {
	ns := s.GetNamespace(name)
	if ns == nil {
		return nil
	}
	return ns.Connections()[id]
}

func (s *server) SetNamespaceConnection(name string, id socket.SocketId, conn *socket.Socket) {
	ns := s.GetNamespace(name)
	if ns == nil {
		return
	}
	ns.Connections()[id] = conn
}

func (s *server) RemoveNamespaceConnection(name string, id socket.SocketId) {
	ns := s.GetNamespace(name)
	if ns == nil {
		return
	}
	delete(ns.Connections(), id)
}

func (s *server) GetNamespaceConnectionsCount(name string) int {
	ns := s.GetNamespace(name)
	if ns == nil {
		return 0
	}
	return len(ns.Connections())
}

func (s *server) RemoveNamespace(name string) {
	for i, ns := range s.namespaces {
		if ns.Name() == name {
			s.namespaces = append(s.namespaces[:i], s.namespaces[i+1:]...)
		}
	}
}

func (s *server) Serve() http.Handler {
	return s.server.ServeHandler(s.opts)
}

func (s *server) Close() {
	s.server.Close(nil)
}
