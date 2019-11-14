package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/let-us-go/zkcli/core"
)

var gitCommit = "unknown"
var built = "unknown"

const version = "0.3.0"

func main() {
	servers := flag.String("s", "127.0.0.1:2181", "Servers")
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	showVersion := flag.Bool("version", false, "Show version info")
	verboseLog := flag.Bool("v", false, "Set to true if want to disable zk log, usefull for use the result in shell script")
	flag.Parse()
	args := flag.Args()

	if *showVersion {
		fmt.Printf("Version:\t%s\nGit commit:\t%s\nBuilt: %s\n",
			version, gitCommit, built)
		os.Exit(0)
	}

	config := core.NewConfig(strings.Split(*servers, ","), !*verboseLog)
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
	defer conn.Close()

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
