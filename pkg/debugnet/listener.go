package debugnet

import "net"

type DebugListener struct {
	li net.Listener
}

func (l *DebugListener) Accept() (net.Conn, error) {
	conn, err := l.li.Accept()
	if err != nil {
		return conn, err
	}

	return NewDebugConn(conn), nil
}

func (l *DebugListener) Close() error {
	return l.li.Close()
}

func (l *DebugListener) Addr() net.Addr {
	return l.li.Addr()
}

func NewDebugListener(li net.Listener) net.Listener {
	return &DebugListener{li: li}
}
