// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Tamir Duberstein (tamird@gmail.com)

package rpcbench

import (
	"math/rand"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"testing"

	"github.com/valyala/gorpc"
	"google.golang.org/grpc"
)

func init() {
	grpc.EnableTracing = false
	gorpc.RegisterType(&Echo2{})
	gorpc.RegisterType(&EchoRequest{})
	gorpc.RegisterType(&EchoResponse{})
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func benchmarkEcho(b *testing.B, size int, accept func(net.Listener) error,
	setup func(net.Addr), teardown func(), setupParallel func() func(string) string) {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		b.Fatal(err)
	}
	go func() {
		if err := accept(listener); err != nil &&
			!strings.HasSuffix(err.Error(), "use of closed network connection") {
			b.Fatal(err)
		}
	}()

	setup(listener.Addr())

	echoMsg := randString(size)

	b.SetBytes(2 * int64(len(echoMsg)))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		runRequest := setupParallel()
		for pb.Next() {
			if a, e := runRequest(echoMsg), echoMsg; a != e {
				b.Fatalf("expected:\n%q\ngot:\n%q", e, a)
			}
		}
	})

	b.StopTimer()

	teardown()
	if err := listener.Close(); err != nil {
		b.Fatal(err)
	}

}

// json-rpc

type Echo struct{}

func (t *Echo) Echo(args *EchoRequest, reply *EchoResponse) error {
	reply.Msg = args.Msg
	return nil
}

func listenAndServeJSONRPC(listener net.Listener) error {
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(new(Echo)); err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func benchmarkEchoJSONRPC(b *testing.B, size int) {
	var client *rpc.Client
	benchmarkEcho(b, size, listenAndServeJSONRPC,
		func(addr net.Addr) {
			var err error
			conn, err := net.Dial(addr.Network(), addr.String())
			if err != nil {
				b.Fatal(err)
			}
			client = jsonrpc.NewClient(conn)
			if err != nil {
				b.Fatal(err)
			}
		},
		func() {
			if err := client.Close(); err != nil {
				b.Fatal(err)
			}
		},
		func() func(string) string {
			return func(echoMsg string) string {
				args := EchoRequest{Msg: echoMsg}
				var reply EchoResponse
				if err := client.Call("Echo.Echo", &args, &reply); err != nil {
					b.Fatal(err)
				}
				return reply.Msg
			}
		},
	)
}

func BenchmarkJSONRPC_1K(b *testing.B) {
	benchmarkEchoJSONRPC(b, 1<<10)
}

func BenchmarkJSONRPC_64K(b *testing.B) {
	benchmarkEchoJSONRPC(b, 64<<10)
}

// gob-rpc

func listenAndServeGobRPC(listener net.Listener) error {
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(new(Echo)); err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go rpcServer.ServeConn(conn)
	}
}

func benchmarkEchoGobRPC(b *testing.B, size int) {
	var client *rpc.Client
	benchmarkEcho(b, size, listenAndServeGobRPC,
		func(addr net.Addr) {
			var err error
			conn, err := net.Dial(addr.Network(), addr.String())
			if err != nil {
				b.Fatal(err)
			}
			client = rpc.NewClient(conn)
			if err != nil {
				b.Fatal(err)
			}
		},
		func() {
			if err := client.Close(); err != nil {
				b.Fatal(err)
			}
		},
		func() func(string) string {
			return func(echoMsg string) string {
				args := EchoRequest{Msg: echoMsg}
				var reply EchoResponse
				if err := client.Call("Echo.Echo", &args, &reply); err != nil {
					b.Fatal(err)
				}
				return reply.Msg
			}
		},
	)
}

func BenchmarkGobRPC_1K(b *testing.B) {
	benchmarkEchoGobRPC(b, 1<<10)
}

func BenchmarkGobRPC_64K(b *testing.B) {
	benchmarkEchoGobRPC(b, 64<<10)
}

// proto-rpc

func listenAndServeProtoRPC(listener net.Listener) error {
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(new(Echo)); err != nil {
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

func benchmarkEchoProtoRPC(b *testing.B, size int) {
	var client *rpc.Client
	benchmarkEcho(b, size, listenAndServeProtoRPC,
		func(addr net.Addr) {
			conn, err := net.Dial(addr.Network(), addr.String())
			if err != nil {
				b.Fatal(err)
			}
			client = NewClient(conn)
		},
		func() {
			if err := client.Close(); err != nil {
				b.Fatal(err)
			}
		},
		func() func(string) string {
			return func(echoMsg string) string {
				args := EchoRequest{Msg: echoMsg}
				var reply EchoResponse
				if err := client.Call("Echo.Echo", &args, &reply); err != nil {
					b.Error(err)
				}
				return reply.Msg
			}
		},
	)
}

func BenchmarkProtoRPC_1K(b *testing.B) {
	benchmarkEchoProtoRPC(b, 1<<10)
}

func BenchmarkProtoRPC_64K(b *testing.B) {
	benchmarkEchoProtoRPC(b, 64<<10)
}

// gorpc

type Echo2 struct{}

func (t *Echo2) Echo(args *EchoRequest) (reply *EchoResponse, err error) {
	reply = new(EchoResponse)
	reply.Msg = args.Msg
	return
}

func newDispatcher() (d *gorpc.Dispatcher) {
	d = gorpc.NewDispatcher()
	se := new(Echo2)
	// Register exported service functions
	// d.AddService("MyServer", &args.MyServer2{})
	d.AddFunc("Echo2.Echo", se.Echo)
	return
}

var c int

func listenAndServeGoRPC(listener net.Listener) error {
	if c == 0 {
		c++
		d := newDispatcher()
		s := gorpc.NewTCPServer(":2678", d.NewHandlerFunc())
		defer s.Stop()
		return s.Serve()
	}
	return nil
}

func benchmarkEchoGoRPC(b *testing.B, size int) {
	var c *gorpc.Client
	var client *gorpc.DispatcherClient
	d := newDispatcher()
	benchmarkEcho(b, size, listenAndServeGoRPC,
		func(addr net.Addr) {
			c = gorpc.NewTCPClient("localhost:2678")
			c.Start()

			// Create a client wrapper for calling server functions.
			client = d.NewFuncClient(c)
		},
		func() {
			c.Stop()
		},
		func() func(string) string {
			return func(echoMsg string) string {
				args := EchoRequest{Msg: echoMsg}
				var reply *EchoResponse
				if b1, err := client.Call("Echo2.Echo", &args); err != nil {
					b.Fatal(err)
				} else if v, ok := b1.(*EchoResponse); !ok {
					b.Fatal("Realy")
				} else {
					reply = v
				}
				return reply.Msg
			}
		},
	)
}

func BenchmarkGoRPC_1K(b *testing.B) {
	benchmarkEchoGoRPC(b, 1<<10)
}

func BenchmarkGoRPC_64K(b *testing.B) {
	benchmarkEchoGoRPC(b, 64<<10)
}
