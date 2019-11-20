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
// Author: Peter Mattis (peter@cockroachlabs.com)

// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpcbench

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/rpc"

	"github.com/gogo/protobuf/proto"
)

type clientCodec struct {
	baseConn

	// temporary work space
	reqBodyBuf   bytes.Buffer
	reqHeaderBuf bytes.Buffer
	reqHeader    RequestHeader
	respHeader   ResponseHeader
}

// NewClientCodec returns a new rpc.ClientCodec using Protobuf-RPC on conn.
func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{
		baseConn: baseConn{
			r: bufio.NewReader(conn),
			w: bufio.NewWriter(conn),
			c: conn,
		},
	}
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) (err error) {
	var request proto.Message
	if param != nil {
		var ok bool
		if request, ok = param.(proto.Message); !ok {
			return fmt.Errorf("protorpc.ClientCodec.WriteRequest: %T does not implement proto.Message", param)
		}
	}

	if err = c.writeRequest(r, request); err != nil {
		return
	}
	return c.w.Flush()
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) (err error) {
	if err = c.readResponseHeader(&c.respHeader); err != nil {
		return
	}

	r.Seq = c.respHeader.Id
	r.ServiceMethod = c.respHeader.Method
	r.Error = c.respHeader.Error
	return
}

func (c *clientCodec) ReadResponseBody(x interface{}) (err error) {
	var response proto.Message
	if x != nil {
		var ok bool
		response, ok = x.(proto.Message)
		if !ok {
			return fmt.Errorf("protorpc.ClientCodec.ReadResponseBody: %T does not implement proto.Message", x)
		}
	}

	if err = c.readResponseBody(&c.respHeader, response); err != nil {
		return
	}

	c.respHeader.Reset()
	return
}

func (c *clientCodec) writeRequest(r *rpc.Request, request proto.Message) (err error) {
	// marshal request
	var pbRequest []byte
	if pbRequest, err = marshal(&c.reqBodyBuf, request); err != nil {
		return
	}

	// generate header
	header := &c.reqHeader
	*header = RequestHeader{
		Id:     r.Seq,
		Method: r.ServiceMethod,
	}

	// marshal header
	var pbHeader []byte
	if pbHeader, err = marshal(&c.reqHeaderBuf, header); err != nil {
		return
	}

	// send header (more)
	if err = c.sendFrame(pbHeader); err != nil {
		return
	}

	return c.sendFrame(pbRequest)
}

func (c *clientCodec) readResponseHeader(header *ResponseHeader) error {
	return c.recvProto(header)
}

func (c *clientCodec) readResponseBody(header *ResponseHeader,
	response proto.Message) error {
	return c.recvProto(response)
}

// NewClient returns a new rpc.Client to handle requests to the
// set of services at the other end of the connection.
func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn))
}
