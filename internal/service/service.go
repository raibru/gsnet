package service

import (
	"fmt"
	"net"
	"os"
	"time"

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
	ctx.Log().Infof("::: use package 'service' wide logging with context: %s", l.contextName)
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
}

// ClientServiceData holds connection data about client services
type ClientServiceData struct {
	Name string
	Addr string
	Port string
}

// GsPktServiceData holds data about groundstation package services
type GsPktServiceData struct {
	Name string
}

// ApplyConnection accept a connection from client and handle incoming data stream
func (s *ServerServiceData) ApplyConnection() error {
	ctx.Log().Infof("::: apply server connection for service %s", s.Name)
	ctx.Log().Tracef("::: ::: create TCP server listener for service %s", s.Name)
	ctx.Log().Tracef("::: ::: listen at %s:%s", s.Addr, s.Port)
	lsn, err := CreateTCPServerListener(s)
	if err != nil {
		ctx.Log().Errorf("::: ::: can not apply server listener due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create listener: %s\n", err.Error())
		return err
	}

	ctx.Log().Tracef("::: ::: establish listener for service %s@%s:%s", s.Name, s.Addr, s.Port)
	for {
		ctx.Log().Trace("::: ::: wait for input...")
		conn, err := lsn.Accept()
		ctx.Log().Trace("::: ::: accept input...")
		if err != nil {
			continue
		}
		go handleServerConnection(conn)
	}
}

// ApplyConnection create a connection to server and handle outgoing data stream
func (s *ClientServiceData) ApplyConnection() error {
	ctx.Log().Infof("::: apply client connection for service %s", s.Name)
	ctx.Log().Tracef("::: ::: create TCP client dialer for service %s", s.Name)
	ctx.Log().Tracef("::: ::: dial to %s:%s", s.Addr, s.Port)
	conn, err := CreateTCPClientConnection(s)
	defer conn.Close()

	if err != nil {
		ctx.Log().Errorf("::: ::: can not apply client dialer due '%s'", err.Error())
		fmt.Fprintf(os.Stderr, "Fatal error create client connection: %s\n", err.Error())
		return err
	}

	for i := range [10]int{} {
		_, err = conn.Write([]byte("HALLO WORLD"))
		if err != nil {
			ctx.Log().Errorf("::: ::: can not write into connection due '%s'", err.Error())
			fmt.Fprintf(os.Stderr, "Error write to connection: %s\n", err.Error())
			return err
		}
		ctx.Log().Tracef("::: ::: successful wrote data into connection %v time(s)", i)
		time.Sleep(5 * time.Second)

	}

	ctx.Log().Infof("::: finish client connection for service %s", s.Name)
	return nil

}

// CreateTCPServerListener create new TCP listener with parameter in ServerService
func CreateTCPServerListener(s *ServerServiceData) (net.Listener, error) {
	serv := s.Addr + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error resolve TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	lsn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error listen TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	return lsn, nil
}

// CreateTCPClientConnection create new TCP connection with parameter in ClientService
func CreateTCPClientConnection(s *ClientServiceData) (net.Conn, error) {
	serv := s.Addr + ":" + s.Port
	addr, err := net.ResolveTCPAddr("tcp4", serv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error resolve dailer TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error connect TCP address %s: %s\n", serv, err.Error())
		return nil, err
	}

	return conn, nil
}

func handleServerConnection(conn net.Conn) {
	ctx.Log().Info("::: ::: handle server connection...")

	defer conn.Close()
	data := make([]byte, 1024)
	for {
		ctx.Log().Trace("::: ::: ::: read data from connection...")
		readLen, err := conn.Read(data)
		if err != nil {
			ctx.Log().Warnf("Failure read data from connetion: %s", err.Error())
			continue
		}

		if readLen == 0 {
			ctx.Log().Info("::: ::: ::: client close connection")
			break // connection already closed by client
		}

		ctx.Log().Tracef("::: ::: ::: successful read data from connection [%s]", data)
		//break

		//doSomething with []byte data
		ctx.Log().Info("::: ::: finish handle server connection")
	}
}

// // https://www.thepolyglotdeveloper.com/2017/05/network-sockets-with-the-go-programming-language/
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
// func startServerMode() { }
//
// func startClientMode() { }
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
