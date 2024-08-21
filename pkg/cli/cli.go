package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const versionFlag = "version"

type RegisterFunc func(fs *pflag.FlagSet)

func AddFlags(v *viper.Viper, cmd *cobra.Command, reg ...RegisterFunc) (*viper.Viper, *cobra.Command) {
	cmd.PersistentFlags().BoolP(versionFlag, "v", false, "Print version and exit")
	for _, r := range reg {
		r(cmd.Flags())
	}
	setupViper(v)
	v.BindPFlags(cmd.Flags())
	return v, cmd
}

func setupViper(v *viper.Viper) {
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
}

func MustRun(run func() error) {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func ShouldPrintVersion(cmd *cobra.Command) bool {
	ok, _ := cmd.PersistentFlags().GetBool(versionFlag)
	return ok
}
