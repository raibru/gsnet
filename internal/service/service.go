package service

import (
	"encoding/hex"
	"net"

	"github.com/raibru/gsnet/internal/archive"
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
var logger = sys.LoggerEntity{}

func (l serviceLogger) Apply() error {
	err := logger.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	logger.Log().Infof("apply service logger behavior: %s", l.contextName)
	logger.Log().Info("::: finish apply service logger")
	return nil
}

func (serviceLogger) Identify() string {
	return logger.ContextName()
}

//
// Services
//

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
	logger.Log().Info("start manage client connections")
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			logger.Log().Info("::: register client connection")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				logger.Log().Info("::: unregister terminated client connection")
			}
		case message := <-manager.broadcast:
			logger.Log().Info("::: broadcast to all managed client connections")
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					logger.Log().Info("::: delete terminated client connections")
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
		logger.Log().Info("::: finish manage client connections")
	}
}

func (manager *ClientManager) receive(client *Client) {
	logger.Log().Info("receive data from managed client connections")
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
			logger.Log().Infof("received data [0x %s]", hexData)

			if manager.service.Archive != nil {
				r := archive.NewRecord(hexData, "RX", "TCP")
				manager.service.Archive <- r
			}
			if manager.service.Forward != nil {
				manager.service.Forward <- data[:length]
			}
			//manager.broadcast <- data
		}
	}
	logger.Log().Info("::: finish receive data")
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	logger.Log().Info("send data to managed client")
	for {
		select {
		case msg, ok := <-client.data:
			if !ok {
				logger.Log().Info("::: finish send data")
				return
			}
			client.socket.Write(msg)
		}
	}
}

func (client *Client) receive() {
	logger.Log().Info("receive data")
	for {
		data := make([]byte, 4096)
		length, err := client.socket.Read(data)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			logger.Log().Infof("::: received data [0x %s]", hex.EncodeToString(data[:length]))
		}
	}
	logger.Log().Info("::: finish receive data")
}

func (client *Client) send() {
	logger.Log().Info("send data")
	for {
		logger.Log().Trace("::: wait for data")
		data := <-client.data

		if string(data) == "EOF" {
			logger.Log().Trace("::: receive EOF flag")
			break
		}

		_, err := client.socket.Write(data)
		if err != nil {
			logger.Log().Errorf("::: failure send data due '%s'", err.Error())
			break
		}
		logger.Log().Trace("::: successful send data")
	}
	logger.Log().Info("::: finish send data")
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
