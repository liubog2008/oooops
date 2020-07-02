package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"

	"github.com/liubog2008/oooops/cmd/operator/app/config"
	"github.com/liubog2008/oooops/cmd/operator/app/options"
	"github.com/liubog2008/oooops/pkg/controller/flow"
	"github.com/liubog2008/oooops/pkg/controller/pipe"
	"github.com/liubog2008/oooops/pkg/version"
)

// NewCommand returns app command
func NewCommand() *cobra.Command {
	opts, err := options.NewOptions()
	if err != nil {
		klog.Fatalf("can't get options: %v", err)
	}
	cmd := &cobra.Command{
		Use:  "operator",
		Long: "operator watches CRD and creates CI/CD flow",
		Run: func(cmd *cobra.Command, args []string) {
			klog.Infof("Version: %v", version.Version())
			printFlags(cmd.Flags())

			cfg, err := opts.Config()
			if err != nil {
				klog.Fatalf("can't parse options to config: %v", err)
			}

			stopCh := make(chan struct{})
			if err := Run(cfg, stopCh); err != nil {
				klog.Fatalf("run operator failed: %v", err)
			}
		},
	}
	opts.AddFlags(cmd.Flags())

	cmd.AddCommand(NewVersionCmd())

	return cmd
}

// NewVersionCmd return cmd reports version
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "version",
		Long: "operator version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %v\n", version.Version())
		},
	}
	return cmd
}

// Run runs the operator
func Run(cfg *config.Config, stopCh chan struct{}) error {
	pc := pipe.NewController(&pipe.ControllerOptions{
		KubeClient: cfg.KubeClient,
		ExtClient:  cfg.ExtClient,

		EventInformer: cfg.EventInformer,
		PipeInformer:  cfg.PipeInformer,
		FlowInformer:  cfg.FlowInformer,
	})

	fc := flow.NewController(&flow.ControllerOptions{
		KubeClient: cfg.KubeClient,
		ExtClient:  cfg.ExtClient,

		FlowInformer:      cfg.FlowInformer,
		JobInformer:       cfg.JobInformer,
		PVCInformer:       cfg.PVCInformer,
		ConfigMapInformer: cfg.ConfigMapInformer,
		PodInformer:       cfg.PodInformer,
	})

	go cfg.KubeInformerFactory.Start(stopCh)
	go cfg.PodInformerFactory.Start(stopCh)
	go cfg.ExtInformerFactory.Start(stopCh)

	go pc.Run(1, stopCh)
	go fc.Run(1, stopCh)

	<-stopCh

	return nil
}

func printFlags(fs *pflag.FlagSet) {
	fs.VisitAll(func(f *pflag.Flag) {
		klog.Infof("FLAG: --%v=%v", f.Name, f.Value)
	})
}
