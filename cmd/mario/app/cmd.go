// Package app defines mario app command
package app

import (
	"fmt"

	"github.com/liubog2008/oooops/cmd/mario/app/config"
	"github.com/liubog2008/oooops/cmd/mario/app/options"
	"github.com/liubog2008/oooops/pkg/mario"
	"github.com/liubog2008/oooops/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

// NewCommand returns app command
func NewCommand() *cobra.Command {
	opts, err := options.NewOptions()
	if err != nil {
		klog.Fatalf("can't get options: %v", err)
	}
	cmd := &cobra.Command{
		Use:  "mario",
		Long: "mario use git to fetch and checkout repo, and serve mario file to others",
		Run: func(cmd *cobra.Command, args []string) {
			klog.Infof("Version: %v", version.Version())
			printFlags(cmd.Flags())

			cfg, err := opts.Config()
			if err != nil {
				klog.Fatalf("can't parse options to config: %v", err)
			}

			stopCh := make(chan struct{})
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
	m, err := mario.New(cfg.RootPath, cfg.RemotePath, cfg.Addr, cfg.Token)
	if err != nil {
		return err
	}
	if err := m.Checkout(cfg.Ref); err != nil {
		return err
	}
	return m.Serve(stopCh)
}

func printFlags(fs *pflag.FlagSet) {
	fs.VisitAll(func(f *pflag.Flag) {
		klog.Infof("FLAG: --%v=%v", f.Name, f.Value)
	})
}
