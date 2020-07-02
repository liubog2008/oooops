package config

import (
	"time"

	"github.com/liubog2008/oooops/pkg/mario/git"
)

type Config struct {
	GitCommand git.Interface

	Addr                    string
	GracefulShutdownTimeout time.Duration

	Remote string
	Ref    string

	Token string
}
