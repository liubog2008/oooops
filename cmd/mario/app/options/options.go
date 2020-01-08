package options

import (
	"os"

	"github.com/liubog2008/oooops/cmd/mario/app/config"
	"github.com/spf13/pflag"
)

// Options defines running options of mario
type Options struct {
	Addr       string
	RootPath   string
	RemotePath string
	Token      string

	Ref string
}

// NewOptions returns new running options
func NewOptions() (*Options, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	opt := &Options{
		Addr:     ":8080",
		RootPath: rootPath,
		Ref:      "master",
	}

	return opt, nil
}

// AddFlags adds flags for mario options
func (opt *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&opt.Addr, "addr", opt.Addr, "listen address")
	fs.StringVar(&opt.RootPath, "root-path", opt.RootPath, "root path of git repo")
	fs.StringVar(&opt.RemotePath, "remote-path", opt.RemotePath, "git repo url")
	fs.StringVar(&opt.Token, "token", opt.Token, "token for download mario file")
	fs.StringVar(&opt.Ref, "ref", opt.Ref, "ref of git repo")
}

// Config parse options to config
func (opt *Options) Config() (*config.Config, error) {
	c := &config.Config{
		RootPath:   opt.RootPath,
		RemotePath: opt.RemotePath,
		Addr:       opt.Addr,
		Token:      opt.Token,
		Ref:        opt.Ref,
	}
	return c, nil
}
