package service

import "github.com/raibru/gsnet/internal/arch"

// PacketServiceData holds connection data about client/server services
type PacketServiceData struct {
	Name     string
	Type     string
	Dialer   *ClientServiceData
	Listener *ServerServiceData
	Archive  *arch.Archive
	Mode     chan string
}

// NewPacketService build new object for listener and dialer service context.
func NewPacketService(
	name string,
	typ string,
	dialer *ClientServiceData,
	listener *ServerServiceData,
	archive *arch.Archive) *PacketServiceData {
	s := &PacketServiceData{
		Name:     name,
		Type:     typ,
		Dialer:   dialer,
		Listener: listener,
		Archive:  archive,
		Mode:     make(chan string),
	}
	return s
}

// ApplyConnection build all dialer/listener connection for current packet service
func (s *PacketServiceData) ApplyConnection() error {
	ctx.Log().Infof("apply all connections for packet service %s", s.Name)
	go func() {
		if err := s.Listener.ApplyConnection(); err != nil {
			ctx.Log().Errorf("Error apply server connection %s: %s", s.Listener.Name, err.Error())
		}
	}()

	go func() {
		if err := s.Dialer.ApplyConnection(); err != nil {
			ctx.Log().Errorf("Error apply dialer connection %s: %s", s.Dialer.Name, err.Error())
		}
	}()

	ctx.Log().Info("::: finish setup channel connections")

	return nil
}
