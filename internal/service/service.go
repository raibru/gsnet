package service

import (
	"encoding/hex"
	"net"
	"time"

	"github.com/raibru/gsnet/internal/arch"
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

			manager.service.Archive.RxCount++
			t := time.Now().Format("2006-01-02 15:04:05.000")
			r := arch.Record{MsgID: manager.service.Archive.RxCount, MsgTime: t, MsgDirection: "RX", Protocol: "TCP", Data: hexData}
			manager.service.Archive.DataChan <- r
			if manager.service.Transfer != nil {
				manager.service.Transfer <- data[:length]
			}
			//manager.broadcast <- data
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
				ctx.Log().Info("::: finish send data")
				return
			}
			client.socket.Write(msg)
		}
	}
}

func (client *Client) receive() {
	ctx.Log().Info("receive data")
	for {
		data := make([]byte, 4096)
		length, err := client.socket.Read(data)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			ctx.Log().Infof("::: received data [0x %s]", hex.EncodeToString(data[:length]))
		}
	}
	ctx.Log().Info("::: finish receive data")
}

func (client *Client) send() {
	ctx.Log().Info("send data")
	for {
		ctx.Log().Trace("::: wait for send data")
		data := <-client.data

		if string(data) == "EOF" {
			ctx.Log().Trace("::: receive EOF flag")
			break
		}

		length, err := client.socket.Write(data)
		if err != nil {
			ctx.Log().Errorf("::: failure send data due '%s'", err.Error())
			break
		}
		ctx.Log().Tracef("::: successful send data: [0x %s]", hex.EncodeToString(data[:length]))
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
