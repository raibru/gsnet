package service

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/raibru/gsnet/internal/archive"
)

// ClientService holds connection data about client services
type ClientService struct {
	Name      string
	Host      string
	Port      string
	Conn      *Client
	archivate chan *archive.Record
	Process   chan []byte // use this chan to acceppt data which have to be processed
	Transfer  chan []byte // use this chan to provide data to transfer somewhere
	Receive   chan []byte // use this chan to handle received data from somewhere
}

// NewClientService deploy a client service with needed data
func NewClientService(name string, host string, port string) *ClientService {
	return &ClientService{
		Name:      name,
		Host:      host,
		Port:      port,
		archivate: nil,
		Process:   nil,
		Transfer:  make(chan []byte),
		Receive:   make(chan []byte),
	}
}

// SetProcess set process data channel
func (s *ClientService) SetProcess(c chan []byte) {
	s.Process = c
}

// SetTransfer set transfer data channel
func (s *ClientService) SetTransfer(c chan []byte) {
	s.Transfer = c
}

// SetReceive set receive data channel
func (s *ClientService) SetReceive(c chan []byte) {
	s.Receive = c
}

// SetArchivate set archive record channel
func (s *ClientService) SetArchivate(c chan *archive.Record) {
	s.archivate = c
}

// ApplyConnection create a connection to server and handle outgoing data stream
func (s *ClientService) ApplyConnection() error {
	logger.Log().Infof("apply client connection for service %s", s.Name)
	logger.Log().Tracef("create TCP client dialer for service %s", s.Name)
	conn, err := CreateTCPClientConnection(s)

	if err != nil {
		logger.Log().Errorf("failure create client TCP connection due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create client TCP connection: %s\n", err.Error())
		return err
	}

	s.Conn = &Client{socket: conn, txData: s.Transfer, rxData: s.Receive}
	go s.Conn.transfer()
	go s.Conn.receive()

	return nil
}

// Finalize cleanup data used by ClientService
func (s *ClientService) Finalize() {
	logger.Log().Infof("finalize service %s", s.Name)
	s.Conn.socket.Close()
	logger.Log().Info("finish finalize service")
}

// ReceivePackets receive from connected server packet data
func (s *ClientService) ReceivePackets(done chan bool) {
	go func() {
		logger.Log().Infof("start service receiving packets from  %s:%s", s.Host, s.Port)

		for {
			data := <-s.Receive
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

// TransferPackets transfers from file readed lines of packets and transfer them
func (s *ClientService) TransferPackets(done chan bool) {
	go func() {
		logger.Log().Infof("start service transfer packets to  %s:%s", s.Host, s.Port)
		for {
			data, more := <-s.Process

			if !more || string(data) == "EOF" {
				logger.Log().Trace("get no more data to transfer notification")
				done <- true
				break
			}

			logger.Log().Tracef("transfer packet: [0x %s]", hex.EncodeToString([]byte(data)))
			s.Conn.txData <- []byte(data)
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
