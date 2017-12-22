package zkcli

import (
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"qiniupkg.com/x/errors.v7"
)

type Config struct {
	Servers  []string
	Username string
	Password string
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

func NewAuth(username, password string) *Auth {
	return &Auth{
		Scheme:  "digest",
		Payload: []byte(fmt.Sprintf("%s:%s", username, password)),
	}
}

func (c *Config) GetAuth() *Auth {
	return NewAuth(c.Username, c.Password)
}

func (c *Config) Connect() (conn *zk.Conn, err error) {
	conn, e, err := zk.Connect(c.Servers, time.Second)
	if err != nil {
		return
	}
	if c.Username != "" && c.Password != "" {
		auth := c.GetAuth()
		conn.AddAuth(auth.Scheme, auth.Payload)
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
