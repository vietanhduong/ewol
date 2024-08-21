package wake

import (
	"net"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	namespace  = "wake"
	ipFlag     = namespace + ".ip"
	portFlag   = namespace + ".port"
	secretFlag = namespace + ".secret"
)

func RegisterFlags(fs *pflag.FlagSet) {
	fs.StringP(ipFlag, "i", "255.255.255.255", "Destination IP address. Unless you have static ARP tables you should use some kind of broadcast address (the broadcast address of the network where the computer resides or the limited broadcast address)")
	fs.Uint16P(portFlag, "p", 9, "Destination port")
	fs.StringP(secretFlag, "s", "", "Secret key which will be used as a simple auth. Only work if you enable serve mode")
}

func InitFromViper(v *viper.Viper, hwaddr net.HardwareAddr) *Wake {
	var ip net.IP
	if raw := v.GetString(ipFlag); raw != "" {
		ip = net.ParseIP(raw)
		if ip == nil {
			log.WithField("ip", raw).Fatal("failed to parse IP address")
		}
	}
	return New(hwaddr, WithIPAddr(ip), WithPort(v.GetUint16(portFlag)), WithSecret(v.GetString(secretFlag)))
}
