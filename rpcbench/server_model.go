package rpcbench

import (
	"net"
	"net/rpc"

	"github.com/cgrates/cgrates/engine"
)

type Responder struct {
	r *engine.Responder
}
func NewResponder(r *engine.Responder)*Responder{
	return &Responder{		r:r	}
}

func (r*Responder)GetCost(arg *CallDescriptor, reply *CallCost)error{
	return r.r.GetCost(arg.Convert(),reply.Convert())
}


func ListenAndServeProtoRPC(address string,r *engine.Responder) error {
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
