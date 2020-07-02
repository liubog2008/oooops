// Package app defines mario app command
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/cmd/mario/app/config"
	"github.com/liubog2008/oooops/cmd/mario/app/options"
	"github.com/liubog2008/oooops/pkg/mario"
	"github.com/liubog2008/oooops/pkg/version"
)

// NewCommand returns app command
func NewCommand() *cobra.Command {
	opts, err := options.NewOptions()
	if err != nil {
		klog.Fatalf("can't get options: %v", err)
	}
	cmd := &cobra.Command{
		Use:  "mario",
		Long: "mario verify working dir, and attach mario onto flow",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := opts.Config()
			if err != nil {
				klog.Fatalf("can't parse options to config: %v", err)
			}

			sig := make(chan os.Signal)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			stopCh := make(chan struct{})

			go func() {
				<-sig
				close(stopCh)
			}()

			if err := Run(cfg, stopCh); err != nil {
				klog.Fatalf("run mario failed: %v", err)
			}
		},
	}

	opts.AddFlags(cmd.Flags())

	cmd.AddCommand(NewVersionCmd())

	return cmd
}

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "version",
		Long: "mario version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %v\n", version.Version())
		},
	}
	return cmd
}

// Run runs the mario
func Run(cfg *config.Config, stopCh chan struct{}) error {
	c := mario.Config{
		GitCommand:              cfg.GitCommand,
		Addr:                    cfg.Addr,
		GracefulShutdownTimeout: cfg.GracefulShutdownTimeout,
		Remote:                  cfg.Remote,
		Ref:                     cfg.Ref,
		Token:                   cfg.Token,
	}
	m := mario.New(&c)

	if err := m.Run(stopCh); err != nil {
		return err
	}
	return nil
}
