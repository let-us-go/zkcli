package main

import (
	"flag"
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
	servers := flag.String("s", "127.0.0.1:2181", "Servers")
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	args := flag.Args()

	conn, e, err := zk.Connect(strings.Split(*servers, ","), time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if *username != "" && *password != "" {
		authS := fmt.Sprintf("%s:%s", *username, *password)
		conn.AddAuth("digest", []byte(authS))
	}
	if len(args) > 0 {
		c := zkcli.ParseCmd(strings.Join(args, " "))
		c.ExitWhenErr = true
		c.Conn = conn
		c.Run()
		return
	}

loop:
	for {
		select {
		case event, ok := <-e:
			if ok && event.State == zk.StateConnected {
				break loop
			}
		}
	}

	p := prompt.New(
		zkcli.GetExecutor(conn),
		zkcli.GetCompleter(conn),
		prompt.OptionTitle("zkcli: interactive zookeeper client"),
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}
