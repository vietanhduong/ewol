package wake

import "net"

type Option func(*Wake)

func WithIPAddr(ip net.IP) Option {
	return func(w *Wake) {
		if ip != nil {
			w.packet.IPAddr = ip
		}
	}
}

func WithPort(port uint16) Option {
	return func(w *Wake) {
		if port > 0 {
			w.packet.Port = port
		}
	}
}

func WithSecret(secret string) Option {
	return func(w *Wake) {
		if secret != "" {
			w.secret = secret
		}
	}
}

func defaultWake() *Wake {
	return &Wake{
		packet: &MagicPacket{
			Port:   MagicPacketPort,
			IPAddr: net.IPv4(255, 255, 255, 255),
		},
	}
}
