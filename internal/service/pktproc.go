package service

import "github.com/raibru/gsnet/internal/archive"

// PacketService holds connection data about client/server services
type PacketService struct {
	Name     string
	Type     string
	Dialer   *ClientService
	Listener *ServerService
	Archive  chan *archive.Record
	Mode     chan string
}

// NewPacketService build new object for listener and dialer service context.
func NewPacketService(name string, typ string) *PacketService {
	return &PacketService{
		Name:     name,
		Type:     typ,
		Dialer:   nil,
		Listener: nil,
		Archive:  nil,
		Mode:     make(chan string),
	}
}

// SetDialer set dialer object
func (s *PacketService) SetDialer(d *ClientService) {
	s.Dialer = d
}

// SetListener set listener object
func (s *PacketService) SetListener(l *ServerService) {
	s.Listener = l
}

// SetArchive set archive record channel
func (s *PacketService) SetArchive(r chan *archive.Record) {
	s.Archive = r
}

// ApplyConnection build all dialer/listener connection for current packet service
func (s *PacketService) ApplyConnection() error {
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

	logger.Log().Info("finish setup channel connections")

	return nil
}
