package rpcbench

import (
	"net"
	"net/rpc"

	"github.com/cgrates/cgrates/engine"
)

type Responder struct {
	r *engine.Responder
}

func NewResponder(r *engine.Responder) *Responder {
	return &Responder{r: r}
}

func (r *Responder) GetCost(arg *CallDescriptor, reply *CallCost) error {
	reply2 := new(engine.CallCost)
	err := r.r.GetCost(arg.Convert(), reply2)
	*reply = *NewCallCost(reply2)
	return err
}

func ListenAndServeProtoRPC(address string, r *engine.Responder) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(NewResponder(r)); err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go ServeConn(rpcServer, conn)
	}
}
