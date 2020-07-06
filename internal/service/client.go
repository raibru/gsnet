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
		process:   nil,
		transfer:  nil,
		receive:   nil,
		//		transfer:  make(chan []byte),
		//		receive:   make(chan []byte),
	}
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
	logger.Log().Infof("apply client connection for service %s", s.Name)
	for i := uint(0); i < s.retry || s.retry == uint(0); i++ {
		logger.Log().Tracef("create TCP client dialer for service %s", s.Name)
		conn, err := CreateTCPClientConnection(s)

		if err != nil {
			logger.Log().Errorf("failure create client TCP connection due '%s'", err.Error())
			fmt.Fprintf(os.Stderr, "failure create client TCP connection: %s\n", err.Error())
			time.Sleep(10 * time.Second)
		} else {
			s.conn = &Client{socket: conn, txData: s.transfer, rxData: s.receive}
			go s.conn.transfer()
			go s.conn.receive()
			return nil
		}
	}
	return errors.New("Failure apply connection")
}

// Finalize cleanup data used by ClientService
func (s *ClientService) Finalize() {
	logger.Log().Infof("finalize service %s", s.Name)
	s.conn.socket.Close()
	logger.Log().Info("finish finalize service")
}

// ReceivePackets receive from connected server packet data
func (s *ClientService) ReceivePackets(done chan bool) {
	go func() {
		logger.Log().Infof("start service receiving packets from  %s:%s", s.Host, s.Port)

		for {
			data := <-s.receive
			logger.Log().Tracef("receive packet: [0x %s]", hex.EncodeToString([]byte(data)))

			if string(data) == "EOF" {
				logger.Log().Trace("get no more data to receive notification")
				done <- true
				break
			}

			if s.archivate != nil {
				hexData := hex.EncodeToString([]byte(data))
				r := archive.NewRecord(hexData, "RX", "TCP")
				s.archivate <- r
			}
		}
		logger.Log().Info("finish receive from client connection")
	}()
}

// PushPackets push packet data via transfer connection
func (s *ClientService) PushPackets(done chan bool) {
	go func() {
		logger.Log().Info("push packet via client service transfer connection")
		for {
			data, more := <-s.process

			if !more || string(data) == "EOF" {
				logger.Log().Trace("get EOF notification from process channel")
				done <- true
				break
			}

			logger.Log().Tracef("transfer packet: [0x %s]", hex.EncodeToString([]byte(data)))
			s.conn.txData <- []byte(data)
			hexData := hex.EncodeToString([]byte(data))

			if s.archivate != nil {
				r := archive.NewRecord(hexData, "TX", "TCP")
				s.archivate <- r
			}
		}

		logger.Log().Info("finish apply client connection")
	}()
}

// CreateTCPClientConnection create new TCP connection with parameter in ClientService
func CreateTCPClientConnection(s *ClientService) (net.Conn, error) {
	logger.Log().Infof("create client dialer service %s with connecting to %s:%s", s.Name, s.Host, s.Port)
	logger.Log().Tracef("resolve tcpTCPAddr %s:%s", s.Host, s.Port)

	serv := s.Host + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		logger.Log().Errorf("failure resolve TCPAddr due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve dailer TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	logger.Log().Tracef("start dial tcp %s@%s:%s", s.Name, s.Host, s.Port)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Log().Errorf("failure connect TCPAddr due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error connect TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	logger.Log().Info("finish create client connection")
	return conn, nil
}
