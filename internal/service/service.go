package service

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/raibru/gsnet/internal/arch"
	"github.com/raibru/gsnet/internal/pkt"
	"github.com/raibru/gsnet/internal/sys"
)

//
// Logging
//

type serviceLogger struct {
	contextName string
}

// LogContext hold logging context
var LogContext = serviceLogger{contextName: "srv"}

// log hold logging context
var ctx = sys.ContextLogger{}

func (l serviceLogger) ApplyLogger() error {
	err := ctx.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	ctx.Log().Infof("apply service logger behavior: %s", l.contextName)
	ctx.Log().Info("::: finish apply service logger")
	return nil
}

func (serviceLogger) GetContextName() string {
	return ctx.ContextName()
}

//
// Services
//

// ServerServiceData holds connection data about server services
type ServerServiceData struct {
	Name string
	Addr string
	Port string
	Arch *arch.Archive
}

// ClientServiceData holds connection data about client services
type ClientServiceData struct {
	Name         string
	Addr         string
	Port         string
	Arch         *arch.Archive
	PacketReader *pkt.InputPacketReader
}

// GsPktServiceData holds data about groundstation package services
type GsPktServiceData struct {
	Name string
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
		//go handleServerConnection(conn)
		ctx.Log().Info("::: finish apply server listener")
	}
}

// ApplyConnection create a connection to server and handle outgoing data stream
func (s *ClientServiceData) ApplyConnection() error {
	ctx.Log().Infof("apply client connection for service %s", s.Name)
	ctx.Log().Tracef("::: create TCP client dialer for service %s", s.Name)
	conn, err := CreateTCPClientConnection(s)
	defer conn.Close()

	if err != nil {
		ctx.Log().Errorf("::: failure create client TCP connection due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create client TCP connection: %s\n", err.Error())
		return err
	}

	client := &Client{socket: conn}
	go client.receive()

	for line := range s.PacketReader.DataChan {
		if line == "EOF" {
			ctx.Log().Trace("::: receive EOF flag")
			break
		}

		hexData := hex.EncodeToString([]byte(line))

		s.Arch.TxCount++
		t := time.Now().Format("2006-01-02 15:04:05.000")
		r := arch.ArchiveRecord{s.Arch.TxCount, t, "TX", "TCP", hexData}
		s.Arch.DataChan <- r

		_, err = conn.Write([]byte(line))
		if err != nil {
			ctx.Log().Errorf("::: failure send data due '%s'", err.Error())
			fmt.Fprintf(os.Stderr, "Error send data: %s\n", err.Error())
			return err
		}
		ctx.Log().Tracef("::: successful send data: [0x %s]", hexData)
		time.Sleep(s.PacketReader.Wait)
	}

	ctx.Log().Info("::: finish apply client connection")
	return nil

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

func handleServerConnection(conn net.Conn) {
	ctx.Log().Info("handle server connection...")

	defer conn.Close()
	data := make([]byte, 4096)
	for {
		ctx.Log().Trace("::: read data from connection...")
		readLen, err := conn.Read(data)
		if err != nil {
			ctx.Log().Errorf("::: failure read data from connetion: %s", err.Error())
			continue
		}

		if readLen == 0 {
			ctx.Log().Info("::: client close connection")
			break // connection already closed by client
		}

		hexData := hex.EncodeToString([]byte(data[:readLen]))
		ctx.Log().Tracef("::: successful read data from connection [%s]", hexData)

		//s.Arch.TxCount++
		//t := time.Now().Format("2006-01-02 15:04:05.000")
		//r := arch.ArchiveRecord{s.Arch.TxCount, t, "TX", "TCP", hexData}
		//s.Arch.DataChan <- r
		//break

		//doSomething with []byte data
		ctx.Log().Info("::: finish handle server connection")
	}
}

// ClientManager hold communication behavior
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	service    *ServerServiceData
}

// Client hold client communication behavior
type Client struct {
	socket net.Conn
	data   chan []byte
}

func (manager *ClientManager) start() {
	ctx.Log().Info("start manage client connections")
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			ctx.Log().Info("::: register client connection")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				ctx.Log().Info("::: unregister terminated client connection")
			}
		case message := <-manager.broadcast:
			ctx.Log().Info("::: broadcast to all managed client connections")
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					ctx.Log().Info("::: delete terminated client connections")
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
		ctx.Log().Info("::: finish manage client connections")
	}
}

func (manager *ClientManager) receive(client *Client) {
	ctx.Log().Info("receive data from managed client connections")
	for {
		data := make([]byte, 4096)
		length, err := client.socket.Read(data)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			hexData := hex.EncodeToString(data[:length])
			ctx.Log().Infof("received data [0x %s]", hexData)

			manager.service.Arch.RxCount++
			t := time.Now().Format("2006-01-02 15:04:05.000")
			r := arch.ArchiveRecord{manager.service.Arch.RxCount, t, "RX", "TCP", hexData}
			manager.service.Arch.DataChan <- r
			//manager.broadcast <- data
		}
	}
	ctx.Log().Info("::: finish receive data")
}

func (client *Client) receive() {
	ctx.Log().Info("receive data")
	for {
		msg := make([]byte, 4096)
		length, err := client.socket.Read(msg)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			ctx.Log().Infof("::: received data [0x %s]", hex.EncodeToString(msg[:length]))
		}
	}
	ctx.Log().Info("::: finish receive data")
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	ctx.Log().Info("send data to managed client")
	for {
		select {
		case msg, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(msg)
		}
	}
	ctx.Log().Info("::: finish send data")
}

// // https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/
//
// // We covered a lot of ground, so it might be easier to look at the application as a whole.
// // Somewhere in your $GOPATH you’ll want a main.go file like previously mentioned. It should
// // contain the following:
//
// package main
//
// import (
//     "bufio"
//     "flag"
//     "fmt"
//     "net"
//     "os"
//     "strings"
// )
//
// type ClientManager struct {
//     clients    map[*Client]bool
//     broadcast  chan []byte
//     register   chan *Client
//     unregister chan *Client
// }
//
// type Client struct {
//     socket net.Conn
//     data   chan []byte
// }
//
// func (manager *ClientManager) start() {
//     for {
//         select {
//         case connection := <-manager.register:
//             manager.clients[connection] = true
//             fmt.Println("Added new connection!")
//         case connection := <-manager.unregister:
//             if _, ok := manager.clients[connection]; ok {
//                 close(connection.data)
//                 delete(manager.clients, connection)
//                 fmt.Println("A connection has terminated!")
//             }
//         case message := <-manager.broadcast:
//             for connection := range manager.clients {
//                 select {
//                 case connection.data <- message:
//                 default:
//                     close(connection.data)
//                     delete(manager.clients, connection)
//                 }
//             }
//         }
//     }
// }
//
// func (manager *ClientManager) receive(client *Client) {
//     for {
//         message := make([]byte, 4096)
//         length, err := client.socket.Read(message)
//         if err != nil {
//             manager.unregister <- client
//             client.socket.Close()
//             break
//         }
//         if length > 0 {
//             fmt.Println("RECEIVED: " + string(message))
//             manager.broadcast <- message
//         }
//     }
// }
//
// func (client *Client) receive() {
//     for {
//         message := make([]byte, 4096)
//         length, err := client.socket.Read(message)
//         if err != nil {
//             client.socket.Close()
//             break
//         }
//         if length > 0 {
//             fmt.Println("RECEIVED: " + string(message))
//         }
//     }
// }
//
// func (manager *ClientManager) send(client *Client) {
//     defer client.socket.Close()
//     for {
//         select {
//         case message, ok := <-client.data:
//             if !ok {
//                 return
//             }
//             client.socket.Write(message)
//         }
//     }
// }
//
// func startServerMode() {
//     fmt.Println("Starting server...")
//     listener, error := net.Listen("tcp", ":12345")
//     if error != nil {
//         fmt.Println(error)
//     }
//     manager := ClientManager{
//         clients:    make(map[*Client]bool),
//         broadcast:  make(chan []byte),
//         register:   make(chan *Client),
//         unregister: make(chan *Client),
//     }
//     go manager.start()
//     for {
//         connection, _ := listener.Accept()
//         if error != nil {
//             fmt.Println(error)
//         }
//         client := &Client{socket: connection, data: make(chan []byte)}
//         manager.register <- client
//         go manager.receive(client)
//         go manager.send(client)
//     }
// }
//
// func startClientMode() {
//     fmt.Println("Starting client...")
//     connection, error := net.Dial("tcp", "localhost:12345")
//     if error != nil {
//         fmt.Println(error)
//     }
//     client := &Client{socket: connection}
//     go client.receive()
//     for {
//         reader := bufio.NewReader(os.Stdin)
//         message, _ := reader.ReadString('\n')
//         connection.Write([]byte(strings.TrimRight(message, "\n")))
//     }
// }
//
// func main() {
//     flagMode := flag.String("mode", "server", "start in client or server mode")
//     flag.Parse()
//     if strings.ToLower(*flagMode) == "server" {
//         startServerMode()
//     } else {
//         startClientMode()
//     }
// }
//
// // If you want to see this application in action, you can execute the following commands:
// //
// // go run *.go --mode server
// // go run *.go --mode client
// // Of course the above two commands should be executed from separate Terminal or Command Prompt windows.
