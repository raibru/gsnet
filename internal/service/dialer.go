package service

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"

	"github.com/raibru/gsnet/internal/archive"
)

// ClientServiceValues holds connection data about client services
type ClientServiceValues struct {
	Name     string
	Host     string
	Port     string
	Conn     *Client
	Process  chan []byte
	Archive  chan *archive.Record
	Transfer chan []byte
}

// NewClientService deploy a client service with needed data
func NewClientService(name string, host string, port string) *ClientServiceValues {
	return &ClientServiceValues{
		Name:     name,
		Host:     host,
		Port:     port,
		Process:  nil,
		Archive:  nil,
		Transfer: make(chan []byte),
	}
}

// SetProcess set process data channel
func (s *ClientServiceValues) SetProcess(p chan []byte) {
	s.Process = p
}

// SetTransfer set transfer data channel
func (s *ClientServiceValues) SetTransfer(t chan []byte) {
	s.Transfer = t
}

// SetArchive set archive record channel
func (s *ClientServiceValues) SetArchive(r chan *archive.Record) {
	s.Archive = r
}

// ApplyConnection create a connection to server and handle outgoing data stream
func (s *ClientServiceValues) ApplyConnection() error {
	logger.Log().Infof("apply client connection for service %s", s.Name)
	logger.Log().Tracef("::: create TCP client dialer for service %s", s.Name)
	conn, err := CreateTCPClientConnection(s)

	if err != nil {
		logger.Log().Errorf("::: failure create client TCP connection due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create client TCP connection: %s\n", err.Error())
		return err
	}

	s.Conn = &Client{socket: conn, data: s.Transfer}
	go s.Conn.send()
	go s.Conn.receive()

	//if s.Process != nil {
	//	go s.Conn.send()
	//}

	//go client.receive()
	//defer func() {
	//	s.Conn.data <- []byte("EOF")
	//}()

	return nil
}

// Finalize cleanup data used by ClientServiceValues
func (s *ClientServiceValues) Finalize() {
	logger.Log().Infof("finalize service %s", s.Name)
	s.Conn.socket.Close()
	//close(s.Archive)
	//close(s.Transfer)
	logger.Log().Info("::: finish finalize service")
}

// SendPackets sends from file readed lines of packets and send them
func (s *ClientServiceValues) SendPackets(done chan bool) {
	go handle(s, done)
}

func handle(s *ClientServiceValues, done chan bool) {
	logger.Log().Infof("send packets from  %s to connected service", s.Name)
	for {
		data, more := <-s.Process

		if !more || string(data) == "EOF" {
			logger.Log().Trace("::: get notify by no more data to send")
			done <- true
			break
		}

		logger.Log().Tracef("::: send packet: [0x %s]", hex.EncodeToString([]byte(data)))
		s.Conn.data <- []byte(data)
		hexData := hex.EncodeToString([]byte(data))

		if s.Archive != nil {
			r := archive.NewRecord(hexData, "TX", "TCP")
			s.Archive <- r
		}
	}

	logger.Log().Info("::: finish apply client connection")
}

// CreateTCPClientConnection create new TCP connection with parameter in ClientService
func CreateTCPClientConnection(s *ClientServiceValues) (net.Conn, error) {
	logger.Log().Infof("create client dialer service %s with connecting to %s:%s", s.Name, s.Host, s.Port)
	logger.Log().Tracef("::: resolve tcpTCPAddr %s:%s", s.Host, s.Port)

	serv := s.Host + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		logger.Log().Errorf("::: failure resolve TCPAddr due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve dailer TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	logger.Log().Tracef("::: start dial tcp %s@%s:%s", s.Name, s.Host, s.Port)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Log().Errorf("::: failure connect TCPAddr due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error connect TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	logger.Log().Info("::: finish create client connection")
	return conn, nil
}
