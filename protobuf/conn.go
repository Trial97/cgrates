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
	"encoding/binary"
	"io"
	"net"

	"github.com/gogo/protobuf/proto"
)

type baseConn struct {
	w        *bufio.Writer
	r        *bufio.Reader
	c        io.Closer
	frameBuf [binary.MaxVarintLen64]byte
}

// Close closes the underlying connection.
func (c *baseConn) Close() error {
	return c.c.Close()
}

func (c *baseConn) sendFrame(data []byte) error {
	// Allocate enough space for the biggest uvarint
	size := c.frameBuf[:]

	if len(data) == 0 {
		n := binary.PutUvarint(size, uint64(0))
		return c.write(c.w, size[:n])
	}

	// Write the size and data
	n := binary.PutUvarint(size, uint64(len(data)))
	if err := c.write(c.w, size[:n]); err != nil {
		return err
	}
	return c.write(c.w, data)
}

func (c *baseConn) write(w io.Writer, data []byte) (err error) {
	for index := 0; index < len(data); {
		var n int
		if n, err = w.Write(data[index:]); err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return
			}
		}
		index += n
	}
	return
}

func (c *baseConn) recvProto(m proto.Message) (err error) {
	var size uint64
	if size, err = binary.ReadUvarint(c.r); err != nil {
		return
	}
	if size == 0 {
		return
	}
	if c.r.Buffered() >= int(size) {
		// Parse proto directly from the buffered data.
		var data []byte
		if data, err = c.r.Peek(int(size)); err != nil {
			return err
		}
		if err = proto.Unmarshal(data, m); err != nil {
			return
		}
		_, err = c.r.Discard(int(size))
		return
	}

	data := make([]byte, size)
	if _, err = io.ReadFull(c.r, data); err != nil {
		return
	}
	return proto.Unmarshal(data, m)
}
