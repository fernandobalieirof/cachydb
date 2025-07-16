package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/fernandobalieirof/cachydb/internal/peer"
)

const defaultListerAddr = ":7481"

var DefaultConfig = Config{
	ListenAddr: defaultListerAddr,
}

type Config struct {
	ListenAddr string
}

type Server struct {
	Config    Config
	peers     map[*peer.Peer]bool
	ln        net.Listener
	addPeerCh chan *peer.Peer
	quitCh    chan struct{}
	msgCh     chan []byte
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListerAddr
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*peer.Peer]bool),
		addPeerCh: make(chan *peer.Peer),
		quitCh:    make(chan struct{}),
		msgCh: make(chan []byte),
	}
}

func (s *Server) Start() error {

	ln, err := net.Listen("tcp", s.Config.ListenAddr)
	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()

	slog.Info("server runnning", "listenAddr", s.Config.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgCh:
			fmt.Println(rawMsg)
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		case <-s.quitCh:
			return

		}

	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {

			slog.Error("failed_to_accept_connection", "error", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := peer.NewPeer(conn, s.msgCh)

	s.addPeerCh <- peer

	slog.Info("new_peer_connected", "remoteAddr", conn.RemoteAddr())

	if err := peer.ReadLoop(); err != nil {
		slog.Error(
			"peer read error",
			"err", err,
			"remoteAddr", conn.RemoteAddr())
	}
}
