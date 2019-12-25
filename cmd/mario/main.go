package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/liubog2008/oooops/cmd/mario/app"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	defer klog.Flush()

	command := app.NewCommand()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
