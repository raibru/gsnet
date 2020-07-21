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
	connType  string
	retry     uint
	conn      *Client
	archivate chan *archive.Record
	process   chan []byte // use this chan to acceppt data which have to be processed
	transfer  chan []byte // use this chan to provide data to transfer somewhere
	receive   chan []byte // use this chan to handle received data from somewhere
}

// NewClientService deploy a client service with needed data
func NewClientService(name string, host string, port string, connType string, retry uint) *ClientService {
	return &ClientService{
		Name:      name,
		Host:      host,
		Port:      port,
		connType:  connType,
		retry:     retry,
		archivate: nil,
		process:   nil,
		transfer:  nil,
		receive:   nil,
	}
}

// SetTransfer set transfer data channel
func (s *ClientService) SetTransfer(c chan []byte) {
	logger.Log().WithField("func", "11102").Trace("... set transfer channel")
	s.transfer = c
	s.conn.txData = c
}

// SetReceive set receive data channel
func (s *ClientService) SetReceive(c chan []byte) {
	logger.Log().WithField("func", "11103").Trace("... set receive channel")
	s.receive = c
	s.conn.rxData = c
}

// SetArchivate set archive record channel
func (s *ClientService) SetArchivate(c chan *archive.Record) {
	logger.Log().WithField("func", "11104").Trace("... set archivate channel")
	s.archivate = c
}

// IsTransferType is current service type a transfer connection
func (s *ClientService) IsTransferType() bool {
	return s.connType == "TX"
}

// IsReceiveType is current service type a receive connection
func (s *ClientService) IsReceiveType() bool {
	return s.connType == "RX"
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
			logger.Log().WithField("func", "11110").Infof("create client socket connection %s", s.Name)
			s.conn = &Client{socket: conn, txData: nil, rxData: nil}
			logger.Log().WithField("func", "11110").Tracef("run successful client socket connection %s", s.Name)
			return nil
		}
	}
	return errors.New("Failure apply client service connection")
}

// Receive receive from connected server packet data
func (s *ClientService) Receive() {
	logger.Log().WithField("func", "11130").Infof("start receive client service  %s", s.Name)

	go func() {
		s.receive = make(chan []byte)
		s.conn.receive(s.receive)

		for {
			logger.Log().WithField("func", "11130").Trace("receive wait incoming data from receive channel")
			select {
			case data, ok := <-s.receive:
				logger.Log().WithField("func", "11130").Tracef("receive data: [0x %s]", hex.EncodeToString([]byte(data)))
				if !ok {
					logger.Log().WithField("func", "11130").Trace("notify by closed receive channel")
					return
				}
				if string(data) == "EOF" {
					logger.Log().WithField("func", "11130").Trace("get EOF notification from receive channel")
					continue
				}

				if s.archivate != nil {
					hexData := hex.EncodeToString([]byte(data))
					r := archive.NewRecord(hexData, "RX", "TCP")
					s.archivate <- r
				}
			}
		}
	}()
}

// Process process incoming packet data and put them into transfer channel
func (s *ClientService) Process(c <-chan []byte) {
	logger.Log().WithField("func", "11150").Infof("start process packets client service  %s", s.Name)
	tc := make(chan []byte)
	s.conn.transfer(tc)

	go func() {
		for {
			logger.Log().WithField("func", "11150").Trace("process packets wait incoming data from process channel")
			select {
			case data, ok := <-c:

				if !ok {
					logger.Log().WithField("func", "11150").Trace("returns from process due closed channel")
					return
				}
				if string(data) == "EOF" {
					logger.Log().WithField("func", "11150").Trace("get EOF notification from process channel")
					break
				}

				//
				// shall do some process stuff with data
				//

				tc <- data
				hexData := hex.EncodeToString([]byte(data))
				if s.archivate != nil {
					r := archive.NewRecord(hexData, "TX", "TCP")
					s.archivate <- r
				}

				//			if s.IsTransferType() {
				//				logger.Log().WithField("func", "11150").Trace("send data into transfer channel")
				//				logger.Log().WithField("func", "11150").Tracef("transfer packet: [0x %s]", hexData)
				//				s.transfer <- []byte(data)
				//
				//				if s.archivate != nil {
				//					r := archive.NewRecord(hexData, "TX", "TCP")
				//					s.archivate <- r
				//				}
				//			} else {
				//				if s.archivate != nil {
				//					r := archive.NewRecord(hexData, "PROC", "INTERNAL")
				//					s.archivate <- r
				//				}
				//				logger.Log().WithField("func", "11150").Trace("finaly data ends in sink")
				//			}
			}
		}
	}()
}

// CreateTCPClientConnection create new TCP connection with parameter in ClientService
func CreateTCPClientConnection(s *ClientService) (net.Conn, error) {
	logger.Log().WithField("func", "11170").Infof("create client dialer service %s connecting %s:%s", s.Name, s.Host, s.Port)
	logger.Log().WithField("func", "11170").Tracef("resolve TCPAddr %s:%s", s.Host, s.Port)

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
