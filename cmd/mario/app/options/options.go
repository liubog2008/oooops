package options

import (
	"os"
	"time"

	"github.com/spf13/pflag"

	"github.com/liubog2008/oooops/cmd/mario/app/config"
	"github.com/liubog2008/oooops/pkg/mario/git"
)

// Options defines running options of mario
type Options struct {
	Remote string
	Ref    string

	Addr                    string
	GracefulShutdownTimeout time.Duration

	Token string
}

// NewOptions returns new running options
func NewOptions() (*Options, error) {
	opt := &Options{
		Addr:                    ":8080",
		GracefulShutdownTimeout: 20 * time.Second,
	}

	return opt, nil
}

// AddFlags adds flags for mario options
func (opt *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&opt.Addr, "addr", opt.Addr, "listen address")

	fs.DurationVar(
		&opt.GracefulShutdownTimeout,
		"graceful-shutdown-timeout",
		opt.GracefulShutdownTimeout,
		"graceful shutdown timeout",
	)

	fs.StringVar(&opt.Remote, "remote", opt.Remote, "remote url of git repo")
	fs.StringVar(&opt.Ref, "ref", opt.Ref, "ref of git repo")

	// TODO(liubog2008): set token which should not be seen in kubernetes yaml
	fs.StringVar(&opt.Token, "token", opt.Token, "token of mario file")
}

// Config parse options to config
func (opt *Options) Config() (*config.Config, error) {
	w, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	gitCmd, err := git.New(w)
	if err != nil {
		return nil, err
	}

	c := &config.Config{
		GitCommand: gitCmd,

		Remote: opt.Remote,
		Ref:    opt.Ref,

		Addr:                    opt.Addr,
		GracefulShutdownTimeout: opt.GracefulShutdownTimeout,

		Token: opt.Token,
	}

	return c, nil
}
