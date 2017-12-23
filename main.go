package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/let-us-go/zkcli/core"
)

const version = "0.1.0"

func main() {
	servers := flag.String("s", "127.0.0.1:2181", "Servers")
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	args := flag.Args()

	config := core.NewConfig(strings.Split(*servers, ","))
	if *username != "" && *password != "" {
		auth := core.NewAuth(
			"digest", fmt.Sprintf("%s:%s", *username, *password),
		)
		config.Auth = auth
	}
	conn, err := config.Connect()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	name, options := core.ParseCmd(strings.Join(args, " "))
	cmd := core.NewCmd(name, options, conn, config)
	if len(args) > 0 {
		cmd.ExitWhenErr = true
		cmd.Run()
		return
	}

	p := prompt.New(
		core.GetExecutor(cmd),
		core.GetCompleter(cmd),
		prompt.OptionTitle("zkcli: A interactive Zookeeper client"),
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}
