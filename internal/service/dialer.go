package service

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/raibru/gsnet/internal/arch"
	"github.com/raibru/gsnet/internal/pkt"
)

// ClientServiceData holds connection data about client services
type ClientServiceData struct {
	Name         string
	Addr         string
	Port         string
	Transfer     chan []byte
	Conn         *Client
	Archive      chan *arch.Record
	PacketReader *pkt.InputPacketReader
}

// NewClientService deploy a client service with needed data
func NewClientService(name string, host string, port string, reader *pkt.InputPacketReader, archSlot chan *arch.Record) *ClientServiceData {
	s := &ClientServiceData{
		Name:         name,
		Addr:         host,
		Port:         port,
		Transfer:     make(chan []byte),
		PacketReader: reader,
		Archive:      archSlot,
	}
	return s
}

// ApplyConnection create a connection to server and handle outgoing data stream
func (s *ClientServiceData) ApplyConnection() error {
	ctx.Log().Infof("apply client connection for service %s", s.Name)
	ctx.Log().Tracef("::: create TCP client dialer for service %s", s.Name)
	conn, err := CreateTCPClientConnection(s)

	if err != nil {
		ctx.Log().Errorf("::: failure create client TCP connection due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create client TCP connection: %s\n", err.Error())
		return err
	}

	//defer conn.Close()

	s.Conn = &Client{socket: conn, data: s.Transfer}
	//go client.receive()
	go s.Conn.send()
	//defer func() {
	//	s.Conn.data <- []byte("EOF")
	//}()

	return nil
}

// Finalize cleanup data used by ClientServiceData
func (s *ClientServiceData) Finalize() {
	ctx.Log().Infof("finalize service %s", s.Name)
	s.Conn.socket.Close()
	close(s.Archive)
	close(s.Transfer)
	ctx.Log().Info("::: finish finalize service")
}

// SendPackets sends from file readed lines of packets and send them
func (s *ClientServiceData) SendPackets() error {
	ctx.Log().Infof("send packets from  %s to connected service", s.Name)
	for line := range s.PacketReader.DataChan {
		ctx.Log().Tracef("::: Send packet: [0x %s]", hex.EncodeToString([]byte(line)))
		if line == "EOF" {
			ctx.Log().Trace("::: read EOF flag from packet reader")
			break
		}

		s.Conn.data <- []byte(line)

		hexData := hex.EncodeToString([]byte(line))

		if s.Archive != nil {
			r := arch.NewRecord(hexData, "TX", "TCP")
			s.Archive <- r
		}

		time.Sleep(s.PacketReader.Wait)
	}

	ctx.Log().Info("::: finish apply client connection")
	return nil

}

// CreateTCPClientConnection create new TCP connection with parameter in ClientService
func CreateTCPClientConnection(s *ClientServiceData) (net.Conn, error) {
	ctx.Log().Infof("create client dialer service %s with connecting to %s:%s", s.Name, s.Addr, s.Port)
	ctx.Log().Tracef("::: resolve tcpTCPAddr %s:%s", s.Addr, s.Port)

	serv := s.Addr + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		ctx.Log().Errorf("::: failure resolve TCPAddr due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve dailer TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	ctx.Log().Tracef("::: start dial tcp %s@%s:%s", s.Name, s.Addr, s.Port)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		ctx.Log().Errorf("::: failure connect TCPAddr due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error connect TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	ctx.Log().Info("::: finish create client connection")
	return conn, nil
}
