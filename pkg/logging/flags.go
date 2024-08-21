package logging

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	namespace  = "log"
	levelFlag  = namespace + ".level"
	formatFlag = namespace + ".format"
)

func RegisterFlags(fs *pflag.FlagSet) {
	fs.String(levelFlag, "info", "Log level")
	fs.String(formatFlag, "text", "Log format. Available options: text, json")
}

func InitFromViper(v *viper.Viper) {
	SetLevel(v.GetString(levelFlag))
	SetFormatter(Formatter(v.GetString(formatFlag)), false)
}
