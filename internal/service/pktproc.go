package service

import "github.com/raibru/gsnet/internal/archive"

// PacketServiceData holds connection data about client/server services
type PacketServiceData struct {
	Name     string
	Type     string
	Dialer   *ClientServiceData
	Listener *ServerServiceData
	Archive  chan *archive.Record
	Mode     chan string
}

// NewPacketService build new object for listener and dialer service context.
func NewPacketService(name string, typ string) *PacketServiceData {
	return &PacketServiceData{
		Name:     name,
		Type:     typ,
		Dialer:   nil,
		Listener: nil,
		Archive:  nil,
		Mode:     make(chan string),
	}
}

// SetDialer set dialer object
func (s *PacketServiceData) SetDialer(d *ClientServiceData) {
	s.Dialer = d
}

// SetListener set listener object
func (s *PacketServiceData) SetListener(l *ServerServiceData) {
	s.Listener = l
}

// SetArchive set archive record channel
func (s *PacketServiceData) SetArchive(r chan *archive.Record) {
	s.Archive = r
}

// ApplyConnection build all dialer/listener connection for current packet service
func (s *PacketServiceData) ApplyConnection() error {
	logger.Log().Infof("apply all connections for packet service %s", s.Name)
	go func() {
		if err := s.Listener.ApplyConnection(); err != nil {
			logger.Log().Errorf("Error apply server connection %s: %s", s.Listener.Name, err.Error())
		}
	}()

	go func() {
		if err := s.Dialer.ApplyConnection(); err != nil {
			logger.Log().Errorf("Error apply dialer connection %s: %s", s.Dialer.Name, err.Error())
		}
	}()

	logger.Log().Info("::: finish setup channel connections")

	return nil
}
