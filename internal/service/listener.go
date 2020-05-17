package service

import (
	"fmt"
	"net"
	"os"

	"github.com/raibru/gsnet/internal/arch"
)

// ServerServiceData holds connection data about server services
type ServerServiceData struct {
	Name    string
	Addr    string
	Port    string
	Linkage *ClientServiceData
	Arch    *arch.Archive
}

// ApplyConnection accept a connection from client and handle incoming data stream
func (s *ServerServiceData) ApplyConnection() error {
	ctx.Log().Infof("apply server connection for service %s", s.Name)
	ctx.Log().Tracef("::: create TCP server listener for service %s", s.Name)
	lsn, err := CreateTCPServerListener(s)
	if err != nil {
		ctx.Log().Errorf("::: failure create TCP server listener due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create TCP listener: %s\n", err.Error())
		return err
	}

	ctx.Log().Tracef("::: establish listener for service %s@%s:%s", s.Name, s.Addr, s.Port)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		service:    s,
	}

	go manager.start()

	for {
		ctx.Log().Trace("::: wait for input...")
		conn, err := lsn.Accept()
		ctx.Log().Trace("::: accept input...")
		if err != nil {
			ctx.Log().Errorf("::: failure accept connection due '%s'", err.Error())
			continue
		}
		client := &Client{socket: conn, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		//go manager.send(client)
		ctx.Log().Info("::: finish apply server listener")
	}
}

// CreateTCPServerListener create new TCP listener with parameter in ServerService
func CreateTCPServerListener(s *ServerServiceData) (net.Listener, error) {
	ctx.Log().Infof("create server listener service %s@%s:%s", s.Name, s.Addr, s.Port)
	ctx.Log().Tracef("::: resolve tcpTCPAddr %s:%s", s.Addr, s.Port)

	serv := s.Addr + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		ctx.Log().Errorf("::: failure resolve TCP address due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve TCP address %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	ctx.Log().Tracef("::: start listen TCP %s@%s:%s", s.Name, s.Addr, s.Port)
	lsn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		ctx.Log().Errorf("::: failure listen TCP due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error listen TCP %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	ctx.Log().Info("::: finish create server listener")
	return lsn, nil
}
