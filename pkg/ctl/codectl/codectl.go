package codectl

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var (
	rootDir string
)

func init() {
	klog.InitFlags(nil)
	cmd.PersistentFlags().StringVarP(&rootDir, "dir", "d", "", "path of the root directory")
}

func readCodeSource(file string) (*v1alpha1.CodeSource, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cs := v1alpha1.CodeSource{}
	if err := yaml.Unmarshal(body, &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}

var cmd = cobra.Command{
	Use:   "codectl",
	Short: "codectl fetch code from SCM by config",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Usage(); err != nil {
			klog.Fatal(err)
		}
	},
}

// Execute executes the command
func Execute() {
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
