package wake

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/vietanhduong/ewol/pkg/logging"
)

var log = logging.WithField("pkg", "network")

const (
	MagicPacketSize = 102
	MagicPacketPort = 9
)

type MagicPacket struct {
	HWAddr net.HardwareAddr
	IPAddr net.IP
	Port   uint16
}

func (p *MagicPacket) Send() error {
	if len(p.IPAddr) == 0 {
		p.IPAddr = net.IPv4(255, 255, 255, 255)
	}

	if p.Port == 0 {
		p.Port = MagicPacketPort
	}

	var packet [MagicPacketSize]byte

	l := log.WithFields(logrus.Fields{
		"remote-ip":   p.IPAddr.String(),
		"remote-port": p.Port,
		"hw-addr":     p.HWAddr.String(),
	})

	copy(packet[0:6], []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	offset := 6

	for i := 0; i < 16; i++ {
		copy(packet[offset:offset+6], p.HWAddr)
		offset += 6
	}

	remote, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.IPAddr, p.Port))
	if err != nil {
		l.WithError(err).Debugf("failed to resolve remote address")
		return fmt.Errorf("failed to resolve remote address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, remote)
	if err != nil {
		l.WithError(err).Debugf("failed to dial remote address")
		return fmt.Errorf("failed to dial remote address: %w", err)
	}

	defer conn.Close()

	if _, err = conn.Write(packet[:]); err != nil {
		l.WithError(err).Debugf("failed to send magic packet")
		return fmt.Errorf("failed to send magic packet: %w", err)
	}

	return nil
}
