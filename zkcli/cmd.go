package zkcli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

const flag int32 = 0

var acl = zk.WorldACL(zk.PermAll)
var ErrUnknownCmd = errors.New("unknown command")

type Cmd struct {
	Name        string
	Options     []string
	ExitWhenErr bool
	Conn        *zk.Conn
	Config      *Config
}

func NewCmd(name string, options []string, conn *zk.Conn, config *Config) *Cmd {
	return &Cmd{
		Name:    name,
		Options: options,
		Conn:    conn,
		Config:  config,
	}
}

func ParseCmd(cmd string) (name string, options []string) {
	args := make([]string, 0)
	for _, cmd := range strings.Split(cmd, " ") {
		if cmd != "" {
			args = append(args, cmd)
		}
	}
	if len(args) == 0 {
		return
	}

	return args[0], args[1:]
}

func (c *Cmd) ls(conn *zk.Conn) (err error) {
	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	children, _, err := conn.Children(p)
	if err != nil {
		return
	}
	fmt.Printf("[%s]\n", strings.Join(children, ", "))
	return
}

func (c *Cmd) get(conn *zk.Conn) (err error) {
	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	value, stat, err := conn.Get(p)
	if err != nil {
		return
	}
	fmt.Printf("%+v\n%s\n", string(value), fmtStat(stat))
	return
}

func (c *Cmd) create(conn *zk.Conn) (err error) {
	p := "/"
	data := ""
	options := c.Options
	if len(options) > 0 {
		p = options[0]
		if len(options) > 1 {
			data = options[1]
		}
	}
	_, err = conn.Create(p, []byte(data), flag, acl)
	if err != nil {
		return
	}
	fmt.Printf("Created %s\n", p)
	root, _ := splitPath(p)
	suggestCache.del(root)
	return
}

func (c *Cmd) set(conn *zk.Conn) (err error) {
	p := "/"
	data := ""
	options := c.Options
	if len(options) > 0 {
		p = options[0]
		if len(options) > 1 {
			data = options[1]
		}
	}
	stat, err := conn.Set(p, []byte(data), -1)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", fmtStat(stat))
	return
}

func (c *Cmd) delete(conn *zk.Conn) (err error) {
	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	err = conn.Delete(p, -1)
	if err != nil {
		return
	}
	fmt.Printf("Deleted %s\n", p)
	root, _ := splitPath(p)
	suggestCache.del(root)
	return
}

func (c *Cmd) connected() bool {
	state := c.Conn.State()
	return state == zk.StateConnected
}

func (c *Cmd) run() (err error) {
	switch c.Name {
	case "ls":
		return c.ls(c.Conn)
	case "get":
		return c.get(c.Conn)
	case "create":
		return c.create(c.Conn)
	case "set":
		return c.set(c.Conn)
	case "delete":
		return c.delete(c.Conn)
	case "close":
		c.Conn.Close()
		if !c.connected() {
			fmt.Println("Closed")
		}
		return
	case "connect":
		conn, err := c.Config.Connect()
		if err != nil {
			return err
		}
		c.Conn = conn
		fmt.Println("Connected")
		return err
	default:
		return ErrUnknownCmd
	}
}

func (c *Cmd) Run() {
	err := c.run()
	if err != nil {
		if err == ErrUnknownCmd {
			printHelp()
			if c.ExitWhenErr {
				os.Exit(2)
			}
		} else {
			printRunError(err)
			if c.ExitWhenErr {
				os.Exit(3)
			}
		}
	}
}

func printHelp() {
	fmt.Println(`get path
ls path
create path data
set path data
delete path
quit
close
connect host:port
addauth scheme auth`)
}

func printRunError(err error) {
	fmt.Println(err)
}

func GetExecutor(cmd *Cmd) func(s string) {
	return func(s string) {
		name, options := ParseCmd(s)
		cmd.Name = name
		cmd.Options = options
		if name == "quit" || name == "exit" {
			os.Exit(0)
		}
		cmd.Run()
	}
}
