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
	connType  string
	archivate chan *archive.Record
	push      chan []byte // use this chan to push data to connection
	process   chan []byte // use this chan to accept data which have to be processed
	forward   chan []byte // use this chan to forward data to somewhere
	notify    chan []byte // use this chan to notify registered clients
}

// NewServerService build new object for listener service context.
// If transfer channel is nil this object is a data sink
func NewServerService(name string, host string, port string, connType string) *ServerService {
	return &ServerService{
		Name:      name,
		Host:      host,
		Port:      port,
		connType:  connType,
		archivate: nil,
		push:      nil,
		process:   nil,
		forward:   nil,
		notify:    nil,
	}
}

// SetPush set push data channel
func (s *ServerService) SetPush(c chan []byte) {
	logger.Log().WithField("func", "11201").Trace("... set push channel")
	s.push = c
}

// SetProcess set process data channel
func (s *ServerService) SetProcess(c chan []byte) {
	logger.Log().WithField("func", "11202").Trace("... set process channel")
	s.process = c
}

// SetForward set forward data channel
func (s *ServerService) SetForward(c chan []byte) {
	logger.Log().WithField("func", "11203").Trace("... set forward channel")
	s.forward = c
}

// SetNotify set notify data channel
func (s *ServerService) SetNotify(c chan []byte) {
	logger.Log().WithField("func", "11204").Trace("... set notify channel")
	s.notify = c
}

// SetArchivate set archive record channel
func (s *ServerService) SetArchivate(r chan *archive.Record) {
	logger.Log().WithField("func", "11205").Trace("... set archivate channel")
	s.archivate = r
}

// IsServiceTransfer is current service type a transfer connection
func (s *ServerService) IsServiceTransfer() bool {
	return s.connType == "TX"
}

// IsServiceReceive is current service type a receive connection
func (s *ServerService) IsServiceReceive() bool {
	return s.connType == "RX"
}

// ApplyConnection accept a connection from client and handle incoming data stream
func (s *ServerService) ApplyConnection() error {
	logger.Log().WithField("func", "11210").Infof("apply server connection for service %s", s.Name)
	logger.Log().WithField("func", "11210").Tracef("create TCP server listener for service %s", s.Name)
	lsn, err := CreateTCPServerListener(s)
	if err != nil {
		logger.Log().WithField("func", "11210").Errorf("failure create TCP server listener due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create TCP listener: %s\n", err.Error())
		return err
	}

	logger.Log().WithField("func", "11210").Tracef("establish listener for service %s@%s:%s", s.Name, s.Host, s.Port)

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		process:    make(chan []byte),
		notify:     make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	s.process = manager.process
	s.notify = manager.notify

	go manager.start()
	go func() {
		for {
			logger.Log().WithField("func", "11210").Trace("apply connection wait for incoming connection")
			conn, err := lsn.Accept()
			logger.Log().WithField("func", "11210").Trace("accept connection")
			if err != nil {
				logger.Log().Errorf("failure accept connection due '%s'", err.Error())
				continue
			}

			logger.Log().WithField("func", "11210").Trace("register new client connection")
			client := &Client{socket: conn, txData: make(chan []byte), rxData: manager.process}
			manager.register <- client

			if s.IsServiceTransfer() {
				logger.Log().WithField("func", "11210").Trace("run client transfer connection")
				//go manager.transfer(client)
				go client.transfer()
			} else if s.IsServiceReceive() {
				logger.Log().WithField("func", "11210").Trace("run client receive connection")
				//go manager.receive(client)
				go client.receive()
			} else {
				logger.Log().Warnf("connection type is not supported in server service:  '%s'", s.connType)
			}

		}
	}()
	logger.Log().WithField("func", "11210").Info("finish apply server listener")
	return nil
}

// PushPackets push packet data via transfer connection
func (s *ServerService) PushPackets(done chan bool) {
	logger.Log().WithField("func", "11220").Info("start push packet to transfer connection")
	for {
		logger.Log().WithField("func", "11220").Trace("push packets wait incoming data from push channel")
		data, more := <-s.push

		if !more || string(data) == "EOF" {
			logger.Log().WithField("func", "11220").Trace("get EOF notification from process channel")
			done <- true
			break
		}

		logger.Log().WithField("func", "11220").Tracef("transfer packet: [0x %s]", hex.EncodeToString([]byte(data)))
		s.notify <- []byte(data)

		if s.archivate != nil {
			hexData := hex.EncodeToString(data)
			r := archive.NewRecord(hexData, "TX", "TCP")
			s.archivate <- r
		}
	}

	logger.Log().WithField("func", "11220").Info("finish push packet to transfer connection")
}

// ProcessPackets processes received packets and transfer them
func (s *ServerService) ProcessPackets() {
	logger.Log().WithField("func", "11230").Info("start process packets service")
	for {
		logger.Log().WithField("func", "11230").Trace("process packets wait incoming data from process channel")
		select {
		case data := <-s.process:
			hexData := hex.EncodeToString(data)
			if s.IsServiceReceive() && s.forward != nil {
				logger.Log().WithField("func", "11230").Trace("put data into forward channel")
				s.forward <- data
				if s.archivate != nil {
					r := archive.NewRecord(hexData, "PROC", "INTERN")
					s.archivate <- r
				}
			} else if s.IsServiceTransfer() && s.notify != nil {
				logger.Log().WithField("func", "11230").Trace("put data into notify channel")
				s.notify <- data
				if s.archivate != nil {
					r := archive.NewRecord(hexData, "TX", "TCP")
					s.archivate <- r
				}
			} else {
				logger.Log().WithField("func", "11230").Trace("sink incomming data")
				if s.archivate != nil {
					r := archive.NewRecord(hexData, "RX", "TCP")
					s.archivate <- r
				}
			}
		}
	}
}

// CreateTCPServerListener create new TCP listener with parameter in ServerService
func CreateTCPServerListener(s *ServerService) (net.Listener, error) {
	logger.Log().WithField("func", "11240").Infof("create server listener service %s@%s:%s", s.Name, s.Host, s.Port)
	logger.Log().WithField("func", "11240").Tracef("resolve tcpTCPAddr %s:%s", s.Host, s.Port)

	serv := s.Host + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		logger.Log().WithField("func", "11240").Errorf("failure resolve TCP address due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve TCP address %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	logger.Log().WithField("func", "11240").Tracef("start listen TCP %s@%s:%s", s.Name, s.Host, s.Port)
	lsn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Log().WithField("func", "11240").Errorf("failure listen TCP due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error listen TCP %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	logger.Log().WithField("func", "11240").Info("finish create server listener")
	return lsn, nil
}
