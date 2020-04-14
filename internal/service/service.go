package service

import (
	"fmt"
	"net"
	"os"
	"time"
)

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

// ApplyTCPService accept a connection from client and handle incoming data stream
func (s *ServerServiceData) ApplyTCPService() error {
	lsn, err := CreateTCPServerListener(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error create listener: %s\n", err.Error())
		return err
	}

	for {
		fmt.Fprintf(os.Stdout, "Server service %s wait for input...\n", s.Name)
		conn, err := lsn.Accept()
		fmt.Fprintf(os.Stdout, "Server service %s accept input...\n", s.Name)
		if err != nil {
			continue
		}
		go handleServerConnection(conn)
	}
}

// ApplyTCPService create a connection to server and handle outgoing data stream
func (s *ClientServiceData) ApplyTCPService() error {
	conn, err := CreateTCPClientConnection(s)
	defer conn.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error create client connection: %s\n", err.Error())
		return err
	}

	for i := range [10]int{} {
		_, err = conn.Write([]byte("HALLO WORLD"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error write to connection: %s\n", err.Error())
			return err
		}
		fmt.Fprintf(os.Stdout, "::write to connection (%v)\n", i)
		time.Sleep(5 * time.Second)

	}

	fmt.Fprintf(os.Stdout, "Succesdul write data into connection\n")
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
	defer conn.Close()
	data := make([]byte, 1024)
	for {
		fmt.Fprintf(os.Stdout, "::service read data...\n")
		readLen, err := conn.Read(data)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Failure read data from client: %s\n", err.Error())
			continue
		}

		if readLen == 0 {
			fmt.Fprintf(os.Stdout, "Client close connection\n")
			break // connection already closed by client
		}

		fmt.Fprintf(os.Stdout, "Succesful read data from client: [%s]\n", data)
		//break

		//doSomething with []byte data
	}
}
