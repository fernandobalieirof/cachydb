package peer

import (
	"net"
)

type Peer struct {
	conn net.Conn
	// sorry future me
	msgCh chan []byte
}

func NewPeer(conn net.Conn, msgCh chan []byte) *Peer {

	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (p *Peer) ReadLoop() error {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}

		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgCh <- msgBuf
	}

}
