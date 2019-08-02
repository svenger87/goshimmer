package tcp

import (
	"net"
	"strconv"
	"sync"

	"github.com/iotaledger/goshimmer/packages/events"
	"github.com/iotaledger/goshimmer/packages/network"
	"github.com/iotaledger/goshimmer/plugins/workerpool"
)

type Server struct {
	Socket net.Listener
	Events serverEvents
	mutex  sync.RWMutex
}

func (this *Server) Shutdown() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.Socket != nil {
		socket := this.Socket
		this.Socket = nil

		socket.Close()
	}
}

func (this *Server) GetSocket() net.Listener {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.Socket
}

func (this *Server) Listen(port int) *Server {
	socket, err := net.Listen("tcp4", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		this.Events.Error.Trigger(err)

		return this
	} else {
		this.Socket = socket
	}

	this.Events.Start.Trigger()
	defer this.Events.Shutdown.Trigger()

	for this.GetSocket() != nil {
		s := this.GetSocket()
		if socket, err := s.Accept(); err != nil {
			if this.GetSocket() != nil {
				this.Events.Error.Trigger(err)
			}
		} else {
			peer := network.NewManagedConnection(socket)
			//TODO: check this
			workerpool.WP.Submit(func() { this.Events.Connect.Trigger(peer) })
		}
	}

	return this
}

func NewServer() *Server {
	return &Server{
		Events: serverEvents{
			Start:    events.NewEvent(events.CallbackCaller),
			Shutdown: events.NewEvent(events.CallbackCaller),
			Connect:  events.NewEvent(managedConnectionCaller),
			Error:    events.NewEvent(events.ErrorCaller),
		},
	}
}
