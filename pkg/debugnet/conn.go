package debugnet

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/http2"
)

type Conn struct {
	conn net.Conn
	ew   *EmptyWriter
}

func NewDebugConn(conn net.Conn) net.Conn {
	return &Conn{conn: conn, ew: &EmptyWriter{}}
}

type EmptyWriter struct {
}

func (ew *EmptyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (c *Conn) Read(b []byte) (n int, err error) {
	n, err = c.conn.Read(b)
	if err == nil {
		log.Printf("Read got %d bytes\n", n)

		c.DebugFrame(b[:n])
	} else {
		if errors.Is(err, io.EOF) {
			log.Printf("got EOF\n")
		}
	}
	return n, err
}

func (c *Conn) DebugFrame(b []byte) {
	reader := bytes.NewReader(b)
	framer := http2.NewFramer(c.ew, reader)
	frame, err := framer.ReadFrame()
	if err != nil {
		log.Printf("ReadFrame got error, %v", err)
	} else {
		log.Printf("%+v", frame)
		if frame.Header().Length == 0 {
			return
		}

		b := make([]byte, frame.Header().Length)
		_n, err := reader.Read(b)
		if err != nil {
			log.Printf("reader.Read got error, %v", err)
		} else {
			log.Printf("content: %s", b[:_n])
		}
	}
}

func (c *Conn) Write(b []byte) (n int, err error) {
	n, err = c.conn.Write(b)
	if err == nil {
		log.Printf("Write %d bytes\n", n)
		c.DebugFrame(b[:n])
	} else {
		if errors.Is(err, io.EOF) {
			log.Printf("got EOF\n")
		}
	}
	return n, err
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
