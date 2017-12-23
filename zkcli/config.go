package zkcli

import (
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"qiniupkg.com/x/errors.v7"
)

type Config struct {
	Servers []string
	Auth    *Auth
}

func NewConfig(Servers []string) *Config {
	return &Config{
		Servers: Servers,
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

func (c *Config) Connect() (conn *zk.Conn, err error) {
	conn, e, err := zk.Connect(c.Servers, time.Second)
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
