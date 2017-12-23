package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/let-us-go/zkcli/zkcli"
)

const version = "0.1.0"

func main() {
	servers := flag.String("s", "127.0.0.1:2181", "Servers")
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Parse()
	args := flag.Args()

	config := zkcli.NewConfig(strings.Split(*servers, ","))
	if *username != "" && *password != "" {
		auth := zkcli.NewAuth(
			"digest", fmt.Sprintf("%s:%s", *username, *password),
		)
		config.Auth = auth
	}
	conn, err := config.Connect()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	name, options := zkcli.ParseCmd(strings.Join(args, " "))
	cmd := zkcli.NewCmd(name, options, conn, config)
	if len(args) > 0 {
		cmd.ExitWhenErr = true
		cmd.Run()
		return
	}

	p := prompt.New(
		zkcli.GetExecutor(cmd),
		zkcli.GetCompleter(cmd),
		prompt.OptionTitle("zkcli: interactive zookeeper client"),
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}
