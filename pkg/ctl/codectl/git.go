package codectl

import (
	"github.com/liubog2008/oooops/pkg/source"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func init() {
	cmd.AddCommand(&gitCmd)
}

var gitCmd = cobra.Command{
	Use:   "git SOURCE",
	Short: "codectl git fetch code from SCM by config",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			if err := cmd.Usage(); err != nil {
				klog.Fatal(err)
			}
		}
		cs, err := readCodeSource(args[0])
		if err != nil {
			klog.Fatalf("can't parse code source: %v", err)
		}
		m, err := source.New(rootDir)
		if err != nil {
			klog.Fatalf("can't create source manager with root dir %s: %v", rootDir, err)
		}
		for _, c := range cs {
			if err := m.Fetch(&c); err != nil {
				klog.Fatalf("can't fetch code: %v", err)
				return
			}
		}
	},
}
