package service

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/raibru/gsnet/internal/archive"
)

// ServerService holds connection data about server services
type ServerService struct {
	Name      string
	Host      string
	Port      string
	archivate chan *archive.Record
	push      chan []byte // use this chan to push data to connection
	process   chan []byte // use this chan to accept data which have to be processed
	forward   chan []byte // use this chan to forward data to somewhere
	notify    chan []byte // use this chan to notify registered clients
}

// NewServerService build new object for listener service context.
// If transfer channel is nil this object is a data sink
func NewServerService(name string, host string, port string) *ServerService {
	return &ServerService{
		Name:      name,
		Host:      host,
		Port:      port,
		archivate: nil,
		push:      nil,
		process:   nil,
		forward:   nil,
		notify:    nil,
	}
}

// SetPush set push data channel
func (s *ServerService) SetPush(c chan []byte) {
	s.push = c
}

// SetProcess set process data channel
func (s *ServerService) SetProcess(c chan []byte) {
	s.process = c
}

// SetForward set forward data channel
func (s *ServerService) SetForward(c chan []byte) {
	s.forward = c
}

// SetNotify set notify data channel
func (s *ServerService) SetNotify(c chan []byte) {
	s.notify = c
}

// SetArchivate set archive record channel
func (s *ServerService) SetArchivate(r chan *archive.Record) {
	s.archivate = r
}

// ApplyConnection accept a connection from client and handle incoming data stream
func (s *ServerService) ApplyConnection() error {
	logger.Log().Infof("apply server connection for service %s", s.Name)
	logger.Log().Tracef("create TCP server listener for service %s", s.Name)
	lsn, err := CreateTCPServerListener(s)
	if err != nil {
		logger.Log().Errorf("failure create TCP server listener due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create TCP listener: %s\n", err.Error())
		return err
	}

	logger.Log().Tracef("establish listener for service %s@%s:%s", s.Name, s.Host, s.Port)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		process:    make(chan []byte),
		notify:     make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	s.process = manager.process
	s.notify = manager.notify

	go s.processPackets()
	go manager.start()
	go func() {
		for {
			logger.Log().Trace("wait for incoming connection")
			conn, err := lsn.Accept()
			logger.Log().Trace("accept connection")
			if err != nil {
				logger.Log().Errorf("failure accept connection due '%s'", err.Error())
				continue
			}

			logger.Log().Trace("register new client connection")
			client := &Client{socket: conn, txData: make(chan []byte), rxData: make(chan []byte)}
			manager.register <- client

			go manager.receive(client)
			go manager.transfer(client)
			go client.receive()
			go client.transfer()

		}
	}()
	logger.Log().Info("finish apply server listener")
	return nil
}

// CreateTCPServerListener create new TCP listener with parameter in ServerService
func CreateTCPServerListener(s *ServerService) (net.Listener, error) {
	logger.Log().Infof("create server listener service %s@%s:%s", s.Name, s.Host, s.Port)
	logger.Log().Tracef("resolve tcpTCPAddr %s:%s", s.Host, s.Port)

	serv := s.Host + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		logger.Log().Errorf("failure resolve TCP address due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve TCP address %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	logger.Log().Tracef("start listen TCP %s@%s:%s", s.Name, s.Host, s.Port)
	lsn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Log().Errorf("failure listen TCP due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error listen TCP %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	logger.Log().Info("finish create server listener")
	return lsn, nil
}

// PushPackets push packet data via transfer connection
func (s *ServerService) PushPackets(done chan bool) {
	go func() {
		logger.Log().Info("start push packet to transfer connection")
		for {
			data, more := <-s.push

			if !more || string(data) == "EOF" {
				logger.Log().Trace("get EOF notification from process channel")
				done <- true
				break
			}

			logger.Log().Tracef("transfer packet: [0x %s]", hex.EncodeToString([]byte(data)))
			s.notify <- []byte(data)

			if s.archivate != nil {
				hexData := hex.EncodeToString(data)
				r := archive.NewRecord(hexData, "TX", "TCP")
				s.archivate <- r
			}
		}

		logger.Log().Info("finish push packet to transfer connection")
	}()
}

// ProcessPackets processes received packets and transfer them
func (s *ServerService) processPackets() {
	logger.Log().Info("start process packets")
	for {
		data, more := <-s.process
		if !more {
			logger.Log().Trace("notify not more data from receive channel")
			break
		}
		if s.archivate != nil {
			hexData := hex.EncodeToString(data)
			r := archive.NewRecord(hexData, "RX", "TCP")
			s.archivate <- r
		}
		if s.forward != nil {
			s.forward <- data
		}
		if s.notify != nil {
			s.notify <- data
		}
	}
	logger.Log().Info("finish process packets")
}
