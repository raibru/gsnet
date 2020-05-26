package service

import (
	"fmt"
	"net"
	"os"

	"github.com/raibru/gsnet/internal/archive"
)

// ServerServiceData holds connection data about server services
type ServerServiceData struct {
	Name     string
	Addr     string
	Port     string
	Transfer chan []byte
	Archive  chan *archive.Record
}

// NewServerService build new object for listener service context.
// If transfer channel is nil this object is a data sink
func NewServerService(name string, host string, port string, transfer chan []byte, archSlot chan *archive.Record) *ServerServiceData {
	s := &ServerServiceData{
		Name:     name,
		Addr:     host,
		Port:     port,
		Transfer: transfer,
		Archive:  archSlot,
	}
	return s
}

// ApplyConnection accept a connection from client and handle incoming data stream
func (s *ServerServiceData) ApplyConnection() error {
	logger.Log().Infof("apply server connection for service %s", s.Name)
	logger.Log().Tracef("::: create TCP server listener for service %s", s.Name)
	lsn, err := CreateTCPServerListener(s)
	if err != nil {
		logger.Log().Errorf("::: failure create TCP server listener due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create TCP listener: %s\n", err.Error())
		return err
	}

	logger.Log().Tracef("::: establish listener for service %s@%s:%s", s.Name, s.Addr, s.Port)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		service:    s,
	}

	go manager.start()

	for {
		logger.Log().Trace("::: wait for input...")
		conn, err := lsn.Accept()
		logger.Log().Trace("::: accept input...")
		if err != nil {
			logger.Log().Errorf("::: failure accept connection due '%s'", err.Error())
			continue
		}
		client := &Client{socket: conn, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		//go manager.send(client)
		logger.Log().Info("::: finish apply server listener")
	}
}

// CreateTCPServerListener create new TCP listener with parameter in ServerService
func CreateTCPServerListener(s *ServerServiceData) (net.Listener, error) {
	logger.Log().Infof("create server listener service %s@%s:%s", s.Name, s.Addr, s.Port)
	logger.Log().Tracef("::: resolve tcpTCPAddr %s:%s", s.Addr, s.Port)

	serv := s.Addr + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		logger.Log().Errorf("::: failure resolve TCP address due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve TCP address %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	logger.Log().Tracef("::: start listen TCP %s@%s:%s", s.Name, s.Addr, s.Port)
	lsn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Log().Errorf("::: failure listen TCP due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error listen TCP %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	logger.Log().Info("::: finish create server listener")
	return lsn, nil
}
