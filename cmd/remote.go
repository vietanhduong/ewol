package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vietanhduong/ewol/pkg/cli"
	"github.com/vietanhduong/ewol/pkg/config"
	"github.com/vietanhduong/ewol/pkg/logging"
)

func newRemoteCmd() *cobra.Command {
	var secret string
	v := viper.New()
	cmd := &cobra.Command{
		Use:   "remote REMOTE_ADDRESS",
		Short: "Remote Wake on LAN for an eWoL server",
		Long:  `NOTE: This command ONLY available for an eWoL server`,
		Args: func(cmd *cobra.Command, args []string) error {
			if cli.ShouldPrintVersion(cmd) {
				return nil
			}
			if len(args) != 1 {
				return fmt.Errorf("remote address is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.InitFromViper(v)
			if cli.ShouldPrintVersion(cmd) {
				config.PrintVersion()
				return nil
			}

			if !strings.HasPrefix(args[0], "http://") && !strings.HasPrefix(args[0], "https://") {
				return errors.New("remote address should start with http:// or https://")
			}

			remote, err := url.Parse(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse remote address: %w", err)
			}

			remote.Path = "/wake"

			req, err := http.NewRequest(http.MethodPost, remote.String(), nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}

			req.Header.Set("User-Agent", config.UserAgent())
			if secret != "" {
				req.Header.Set("Authorization", secret)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to send request: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				msg := fmt.Sprintf("failed to send request: %s", resp.Status)
				b, _ := io.ReadAll(resp.Body)
				if b != nil {
					msg = fmt.Sprintf("%s: %s", msg, string(b))
				}
				return errors.New(msg)
			}
			fmt.Println("Packet sent successfully!!!")
			return nil
		},
	}
	cli.AddFlags(v, cmd, logging.RegisterFlags)
	cmd.Flags().StringVarP(&secret, "wake.secret", "s", "", "Secret key to wake up the device")
	return cmd
}
