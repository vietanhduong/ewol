package wake

import (
	"net"
	"net/http"
)

type Wake struct {
	packet *MagicPacket
	secret string
}

func New(hwaddr net.HardwareAddr, opt ...Option) *Wake {
	w := defaultWake()
	w.packet.HWAddr = hwaddr
	for _, o := range opt {
		o(w)
	}
	w.packet.HWAddr = hwaddr
	return w
}

func (w *Wake) Send() error {
	return w.packet.Send()
}

func (w *Wake) HttpHandler() (string, http.Handler) {
	return "/wake", http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(wr, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Authorization") != w.secret {
			http.Error(wr, "unauthorized", http.StatusUnauthorized)
			return
		}

		if err := w.Send(); err != nil {
			http.Error(wr, "internal server error", http.StatusInternalServerError)
			log.WithError(err).Errorf("failed to send magic packet")
			return
		}

		wr.WriteHeader(http.StatusAccepted)
	})
}

func (w *Wake) IsSecretEmpty() bool { return w.secret == "" }
