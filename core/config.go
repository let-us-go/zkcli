package core

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

type Config struct {
	Servers []string
	Auth    *Auth

	SlientLog bool
}

func NewConfig(Servers []string, slientLog bool) *Config {
	return &Config{
		Servers:   Servers,
		SlientLog: slientLog,
	}
}

type Auth struct {
	Scheme  string
	Payload []byte
}

func NewAuth(scheme, auth string) *Auth {
	return &Auth{
		Scheme:  scheme,
		Payload: []byte(auth),
	}
}

type emptyLogger struct{}

func (emptyLogger) Printf(format string, a ...interface{}) {
	// do nothing
}

func (c *Config) Connect() (conn *zk.Conn, err error) {
	logger := zk.WithLogger(zk.DefaultLogger)
	if c.SlientLog {
		logger = zk.WithLogger(emptyLogger{})
	}
	conn, e, err := zk.Connect(c.Servers, time.Second, logger)
	if err != nil {
		return
	}
	if c.Auth != nil {
		auth := c.Auth
		err = conn.AddAuth(auth.Scheme, auth.Payload)
		if err != nil {
			return
		}
	}
	n := 0
	failed := false
loop:
	for {
		select {
		case event, ok := <-e:
			n += 1
			if ok && event.State == zk.StateConnected {
				break loop
			} else if n > 3 {
				failed = true
				break loop
			}
		}
	}
	if failed {
		err = errors.New(
			fmt.Sprintf("Failed to connect to %s!", strings.Join(c.Servers, ",")),
		)
	}
	return
}
