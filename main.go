package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/let-us-go/zkcli/zkcli"
	"github.com/samuel/go-zookeeper/zk"
)

const version = "0.1.0"

func main() {
	args := os.Args[1:]
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(args) > 0 {
		c := zkcli.ParseCmd(strings.Join(args, " "))
		c.ExitWhenErr = true
		c.Conn = conn
		c.Run()
		return
	}

	p := prompt.New(
		zkcli.GetExecutor(conn),
		zkcli.GetCompleter(conn),
		prompt.OptionTitle("zkcli: interactive zookeeper client"),
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}
