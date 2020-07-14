package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/raibru/gsnet/internal/archive"
)

// ClientService holds connection data about client services
type ClientService struct {
	Name      string
	Host      string
	Port      string
	retry     uint
	conn      *Client
	archivate chan *archive.Record
	push      chan []byte // use this chan to push data to connection
	process   chan []byte // use this chan to acceppt data which have to be processed
	transfer  chan []byte // use this chan to provide data to transfer somewhere
	receive   chan []byte // use this chan to handle received data from somewhere
}

// NewClientService deploy a client service with needed data
func NewClientService(name string, host string, port string, retry uint) *ClientService {
	return &ClientService{
		Name:      name,
		Host:      host,
		Port:      port,
		retry:     retry,
		archivate: nil,
		push:      nil,
		process:   nil,
		transfer:  nil,
		receive:   nil,
		//		transfer:  make(chan []byte),
		//		receive:   make(chan []byte),
	}
}

// SetPush set push data channel
func (s *ClientService) SetPush(c chan []byte) {
	s.process = c
}

// SetProcess set process data channel
func (s *ClientService) SetProcess(c chan []byte) {
	s.process = c
}

// SetTransfer set transfer data channel
func (s *ClientService) SetTransfer(c chan []byte) {
	s.transfer = c
}

// SetReceive set receive data channel
func (s *ClientService) SetReceive(c chan []byte) {
	s.receive = c
}

// SetArchivate set archive record channel
func (s *ClientService) SetArchivate(c chan *archive.Record) {
	s.archivate = c
}

// ApplyConnection create a connection to server and handle outgoing data stream
func (s *ClientService) ApplyConnection() error {
	logger.Log().WithField("func", "11110").Infof("apply client connection for service %s", s.Name)
	for i := uint(0); i < s.retry || s.retry == uint(0); i++ {
		logger.Log().WithField("func", "11110").Tracef("create TCP client dialer for service %s", s.Name)
		conn, err := CreateTCPClientConnection(s)

		if err != nil {
			logger.Log().WithField("func", "11110").Errorf("failure create client TCP connection due '%s'", err.Error())
			fmt.Fprintf(os.Stderr, "failure create client TCP connection: %s\n", err.Error())
			time.Sleep(10 * time.Second)
		} else {
			logger.Log().WithField("func", "11110").Infof("create client connection %s", s.Name)
			ctx := make(chan []byte)
			crx := make(chan []byte)
			s.SetTransfer(ctx)
			s.SetReceive(crx)
			s.conn = &Client{socket: conn, txData: ctx, rxData: crx}
			go s.conn.transfer()
			go s.conn.receive()
			logger.Log().WithField("func", "11110").Tracef("run successful client connection %s", s.Name)
			return nil
		}
	}
	return errors.New("Failure apply connection")
}

// Finalize cleanup data used by ClientService
func (s *ClientService) Finalize() {
	logger.Log().WithField("func", "11120").Infof("finalize service %s", s.Name)
	s.conn.socket.Close()
	logger.Log().WithField("func", "11120").Info("finish finalize service")
}

// ReceivePackets receive from connected server packet data
func (s *ClientService) ReceivePackets() {
	logger.Log().WithField("func", "11130").Infof("start service receiving packets from  %s:%s", s.Host, s.Port)

	for {
		logger.Log().WithField("func", "11130").Trace("receive packets wait incoming data from receive channel")
		select {
		case data := <-s.receive:
			logger.Log().WithField("func", "11130").Tracef("receive packet: [0x %s]", hex.EncodeToString([]byte(data)))

			if string(data) == "EOF" {
				logger.Log().WithField("func", "11130").Trace("get EOF notification from receive channel")
				continue
			}

			if s.archivate != nil {
				hexData := hex.EncodeToString([]byte(data))
				r := archive.NewRecord(hexData, "RX", "TCP")
				s.archivate <- r
			}
			//s.process <- data
		}
	}
}

// NotifyPackets notify incoming data from client service receive channel
func (s *ClientService) NotifyPackets() {
	logger.Log().WithField("func", "11140").Info("start notify server service")

	for {
		logger.Log().WithField("func", "11140").Trace("notify packets wait incoming data from receive channel")
		select {
		case data := <-s.receive:
			logger.Log().WithField("func", "11140").Tracef("notify packet: [0x %s]", hex.EncodeToString([]byte(data)))

			if string(data) == "EOF" {
				logger.Log().WithField("func", "11140").Trace("get EOF notification from receive channel")
				continue
			}

			if s.archivate != nil {
				hexData := hex.EncodeToString([]byte(data))
				r := archive.NewRecord(hexData, "RX", "TCP")
				s.archivate <- r
			}
			logger.Log().WithField("func", "11140").Trace("put data to client service process channel")
			s.process <- data
		}
	}
}

// PushPackets push packet data via transfer connection
func (s *ClientService) PushPackets(done chan bool) {
	logger.Log().WithField("func", "11150").Info("push packet via client service transfer connection")
	for {
		logger.Log().WithField("func", "11150").Trace("push packets wait incoming data from process channel")
		data, more := <-s.process

		if !more || string(data) == "EOF" {
			logger.Log().WithField("func", "11150").Trace("get EOF notification from process channel")
			done <- true
			break
		}

		logger.Log().WithField("func", "11150").Tracef("push packet: [0x %s]", hex.EncodeToString([]byte(data)))
		s.transfer <- []byte(data)
		hexData := hex.EncodeToString([]byte(data))

		if s.archivate != nil {
			r := archive.NewRecord(hexData, "TX", "TCP")
			s.archivate <- r
		}
	}
}

// ProcessPackets process a packet into another state before transfer it
func (s *ClientService) ProcessPackets() {
	logger.Log().WithField("func", "11160").Info("start process packets service")
	for {
		logger.Log().WithField("func", "11160").Trace("process packets wait incoming data from process channel")
		select {
		case data := <-s.process:
			hexData := hex.EncodeToString([]byte(data))
			logger.Log().WithField("func", "11160").Tracef("process packet: [0x %s]", hexData)

			if string(data) == "EOF" {
				logger.Log().WithField("func", "11160").Trace("get EOF notification from receive channel")
				continue
			}

			if s.archivate != nil {
				r := archive.NewRecord(hexData, "PROC", "intern")
				s.archivate <- r
			}
		}
	}
}

// CreateTCPClientConnection create new TCP connection with parameter in ClientService
func CreateTCPClientConnection(s *ClientService) (net.Conn, error) {
	logger.Log().WithField("func", "11170").Infof("create client dialer service %s with connecting to %s:%s", s.Name, s.Host, s.Port)
	logger.Log().WithField("func", "11170").Tracef("resolve tcpTCPAddr %s:%s", s.Host, s.Port)

	serv := s.Host + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		logger.Log().WithField("func", "11170").Errorf("failure resolve TCPAddr due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve dailer TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	logger.Log().WithField("func", "11170").Tracef("start dial tcp %s@%s:%s", s.Name, s.Host, s.Port)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Log().WithField("func", "11170").Errorf("failure connect TCPAddr due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error connect TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	logger.Log().WithField("func", "11170").Info("finish create client connection")
	return conn, nil
}
