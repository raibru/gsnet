package service

import (
	"fmt"
	"net"
	"os"

	"github.com/raibru/gsnet/internal/archive"
)

// ServerServiceData holds connection data about server services
type ServerServiceData struct {
	Name    string
	Host    string
	Port    string
	Process chan []byte
	Forward chan []byte
	Archive chan *archive.Record
}

// NewServerService build new object for listener service context.
// If transfer channel is nil this object is a data sink
func NewServerService(name string, host string, port string) *ServerServiceData {
	return &ServerServiceData{
		Name:    name,
		Host:    host,
		Port:    port,
		Process: nil,
		Forward: nil,
		Archive: nil,
	}
}

// SetProcess set process data channel
func (s *ServerServiceData) SetProcess(p chan []byte) {
	s.Process = p
}

// SetForward set forward data channel
func (s *ServerServiceData) SetForward(f chan []byte) {
	s.Forward = f
}

// SetArchive set archive record channel
func (s *ServerServiceData) SetArchive(r chan *archive.Record) {
	s.Archive = r
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

	logger.Log().Tracef("::: establish listener for service %s@%s:%s", s.Name, s.Host, s.Port)

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
	logger.Log().Infof("create server listener service %s@%s:%s", s.Name, s.Host, s.Port)
	logger.Log().Tracef("::: resolve tcpTCPAddr %s:%s", s.Host, s.Port)

	serv := s.Host + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		logger.Log().Errorf("::: failure resolve TCP address due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error resolve TCP address %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	logger.Log().Tracef("::: start listen TCP %s@%s:%s", s.Name, s.Host, s.Port)
	lsn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Log().Errorf("::: failure listen TCP due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error listen TCP %s: %s\n", s.Name, err.Error())
		return nil, err
	}

	logger.Log().Info("::: finish create server listener")
	return lsn, nil
}

// BroadcastPackets broadcast packet data to managed clients
func (s *ServerServiceData) BroadcastPackets(done chan bool) {
	go broadcast(s, done)
}

func broadcast(s *ServerServiceData, done chan bool) {
	logger.Log().Infof("send packets from  %s to managed client services", s.Name)
	for {
		_ = s
		done <- true
		//		data, more := <-s.Process
		//
		//		if !more || string(data) == "EOF" {
		//			logger.Log().Trace("::: get notify by no more data to send")
		//			done <- true
		//			break
		//		}
		//
		//		logger.Log().Tracef("::: send packet: [0x %s]", hex.EncodeToString([]byte(data)))
		//		s.Conn.data <- []byte(data)
		//		hexData := hex.EncodeToString([]byte(data))
		//
		//		if s.Archive != nil {
		//			r := archive.NewRecord(hexData, "TX", "TCP")
		//			s.Archive <- r
		//		}
	}

	logger.Log().Info("::: finish apply client connection")
}
