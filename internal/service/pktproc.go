package service

import "github.com/raibru/gsnet/internal/archive"

// PacketService holds connection data about client/server services
type PacketService struct {
	Name      string
	Type      string
	dialer    *ClientService
	listener  *ServerService
	archivate chan *archive.Record
	Mode      chan string
}

// NewPacketService build new object for listener and dialer service context.
func NewPacketService(name string, typ string) *PacketService {
	return &PacketService{
		Name:      name,
		Type:      typ,
		dialer:    nil,
		listener:  nil,
		archivate: nil,
		Mode:      make(chan string),
	}
}

// SetDialer set dialer object
func (s *PacketService) SetDialer(d *ClientService) {
	s.dialer = d
}

// SetListener set listener object
func (s *PacketService) SetListener(l *ServerService) {
	s.listener = l
}

// SetArchivate set archive record channel
func (s *PacketService) SetArchivate(r chan *archive.Record) {
	s.archivate = r
}

// ApplyConnection build all dialer/listener connection for current packet service
func (s *PacketService) ApplyConnection() error {
	logger.Log().WithField("func", "11310").Infof("apply all connections for packet service %s", s.Name)

	logger.Log().WithField("func", "11310").Info("call apply listener connection")
	go func() {
		if err := s.listener.ApplyConnection(); err != nil {
			logger.Log().WithField("func", "11310").Errorf("Error apply server connection %s: %s", s.listener.Name, err.Error())
		}
		s.listener.Process()
	}()

	logger.Log().WithField("func", "11310").Info("call apply dialer connection")
	go func() {
		if err := s.dialer.ApplyConnection(); err != nil {
			logger.Log().WithField("func", "11310").Errorf("Error apply dialer connection %s: %s", s.dialer.Name, err.Error())
		}
		s.dialer.Receive()
	}()

	logger.Log().WithField("func", "11310").Info("finish apply connections for packet service")

	return nil
}
