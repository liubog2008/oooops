package imagectl

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/liubog2008/oooops/pkg/apis/flow/v1alpha1"
	"github.com/liubog2008/oooops/pkg/image"
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

func readImages(file string) ([]v1alpha1.Image, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	images := []v1alpha1.Image{}
	if err := yaml.Unmarshal(body, &images); err != nil {
		return nil, err
	}
	return images, nil
}

var cmd = cobra.Command{
	Use:   "imagectl PATH",
	Short: "imagectl builds and pushes images to remote registry",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			if err := cmd.Usage(); err != nil {
				klog.Fatal(err)
			}
		}
		images, err := readImages(args[0])
		if err != nil {
			klog.Fatalf("can't find images list: %v", err)
		}
		m, err := image.New(rootDir)
		if err != nil {
			klog.Fatalf("can't create image manager: %v", err)
		}
		for _, im := range images {
			if err := m.BuildAndPush(&im); err != nil {
				klog.Fatalf("can't build and push image %v", im)
			}
		}
	},
}

// Execute executes the command
func Execute() {
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
