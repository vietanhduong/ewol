package cmd

import (
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vietanhduong/ewol/pkg/cli"
	"github.com/vietanhduong/ewol/pkg/config"
	"github.com/vietanhduong/ewol/pkg/logging"
	"github.com/vietanhduong/ewol/pkg/server"
	"github.com/vietanhduong/ewol/pkg/wake"
)

func New() *cobra.Command {
	var serve bool
	v := viper.New()
	cmd := &cobra.Command{
		Use:   "ewol HARDWARE_ADDRESS",
		Short: "Extended Wake-on-LAN",
		Long: `Extended Wake-on-LAN is a tool to wake up devices on a local network.
You can also publish the service to the network and wake up the input device remotely via an API call.`,
		Example: `# Wake-on-LAN directly
$ ewol 00:11:22:33:44:55

# Publish the service to the network and wake up the device remotely
# You can also provide the wake secret via the WAKE_SECRET environment variable
$ ewol 00:11:22:33:44:55 --serve --wake.secret mysecret

# Specify the IP address and port
$ ewol 00:11:22:33:44:55 --wake.ip 255.255.255.255 --wake.port 9
  `,
		Args: func(cmd *cobra.Command, args []string) error {
			if cli.ShouldPrintVersion(cmd) {
				return nil
			}
			if len(args) != 1 {
				return fmt.Errorf("hardware address is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.InitFromViper(v)

			if cli.ShouldPrintVersion(cmd) {
				config.PrintVersion()
				return nil
			}

			hwaddr, err := net.ParseMAC(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse hardware address: %w", err)
			}

			w := wake.InitFromViper(v, hwaddr)
			if !serve {
				return w.Send()
			}

			if w.IsSecretEmpty() {
				logging.Warn("You are running in serve mode without a secret. This is not recommended.")
			}

			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()

			s := server.InitFromViper(v)
			s.RegisterHandler(w)
			return s.Run(ctx.Done())
		},
	}

	cmd.Flags().BoolVar(&serve, "serve", false, "Enable serve mode. This will create an HTTP server to listen for incoming requests")
	cli.AddFlags(v, cmd, logging.RegisterFlags, server.RegisterFlags, wake.RegisterFlags)
	cmd.AddCommand(newRemoteCmd())
	return cmd
}
