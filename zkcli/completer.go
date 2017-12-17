package zkcli

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/samuel/go-zookeeper/zk"
)

func GetCompleter(conn *zk.Conn) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		if d.TextBeforeCursor() == "" {
			return []prompt.Suggest{}
		}
		args := strings.Split(d.TextBeforeCursor(), " ")
		return argumentsCompleter(excludeOptions(args), conn)
	}
}

var commands = []prompt.Suggest{
	{Text: "get", Description: "Get data"},
	{Text: "ls", Description: "List children"},
	{Text: "create", Description: "Create a node"},
	{Text: "set", Description: "Update a node"},
	{Text: "delete", Description: "Delete a node"},
	{Text: "quit", Description: "Exit this program"},
	{Text: "exit", Description: "Exit this program"},
}

var childrenSuggests = map[string][]prompt.Suggest{}

func argumentsCompleter(args []string, conn *zk.Conn) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}

	first := args[0]
	switch first {
	case "get":
		path := args[1]
		ps := strings.Split(path, "/")
		root := strings.Join(ps[:len(ps)-1], "/")
		if root == "" {
			root = "/"
		}
		return prompt.FilterContains(getChildrenCompletions(conn, root), path, true)
	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}

func getChildrenCompletions(conn *zk.Conn, root string) []prompt.Suggest {
	if value, ok := childrenSuggests[root]; ok {
		return value
	}

	children, _, err := conn.Children(root)
	if err != nil || len(children) == 0 {
		return []prompt.Suggest{}
	}

	s := make([]prompt.Suggest, len(children))
	for i, child := range children {
		path := "/"
		if root == "/" {
			path = fmt.Sprintf("/%s", child)
		} else {
			path = fmt.Sprintf("%s/%s", root, child)
		}
		s[i] = prompt.Suggest{
			Text: path,
		}
	}
	childrenSuggests[root] = s
	return s
}
