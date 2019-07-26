package udp

import (
	"net"
	"strconv"
	"sync"

	"github.com/iotaledger/goshimmer/packages/events"
)

type Server struct {
	Socket            net.PacketConn
	ReceiveBufferSize int
	Events            serverEvents
	mutex             sync.RWMutex
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

func (this *Server) GetSocket() net.PacketConn {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.Socket
}

func (this *Server) Listen(address string, port int) {
	if socket, err := net.ListenPacket("udp", address+":"+strconv.Itoa(port)); err != nil {
		this.Events.Error.Trigger(err)

		return
	} else {
		this.Socket = socket
	}

	this.Events.Start.Trigger()
	defer this.Events.Shutdown.Trigger()

	buf := make([]byte, this.ReceiveBufferSize)
	for this.GetSocket() != nil {
		s := this.GetSocket()
		if bytesRead, addr, err := s.ReadFrom(buf); err != nil {
			if this.GetSocket() != nil {
				this.Events.Error.Trigger(err)
			}
		} else {
			this.Events.ReceiveData.Trigger(addr.(*net.UDPAddr), buf[:bytesRead])
		}
	}
}

func NewServer(receiveBufferSize int) *Server {
	return &Server{
		ReceiveBufferSize: receiveBufferSize,
		Events: serverEvents{
			Start:       events.NewEvent(events.CallbackCaller),
			Shutdown:    events.NewEvent(events.CallbackCaller),
			ReceiveData: events.NewEvent(dataCaller),
			Error:       events.NewEvent(events.ErrorCaller),
		},
	}
}
