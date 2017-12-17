package zkcli

import (
	"fmt"
	"strings"

	"path"

	"github.com/c-bata/go-prompt"
	"github.com/samuel/go-zookeeper/zk"
)

var commands = []prompt.Suggest{
	{Text: "get", Description: "Get data"},
	{Text: "ls", Description: "List children"},
	{Text: "create", Description: "Create a node"},
	{Text: "set", Description: "Update a node"},
	{Text: "delete", Description: "Delete a node"},
	{Text: "quit", Description: "Exit this program"},
	{Text: "exit", Description: "Exit this program"},
}

var suggestCache = newSuggestCache()

func GetCompleter(conn *zk.Conn) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		if d.TextBeforeCursor() == "" {
			return []prompt.Suggest{}
		}
		args := strings.Split(d.TextBeforeCursor(), " ")
		return argumentsCompleter(excludeOptions(args), conn)
	}
}

func argumentsCompleter(args []string, conn *zk.Conn) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}

	first := args[0]
	switch first {
	case "get", "ls", "create", "set", "delete":
		p := args[1]
		if len(args) > 2 {
			return []prompt.Suggest{}
		}
		root, _ := splitPath(p)
		return prompt.FilterContains(getChildrenCompletions(conn, root), p, true)
	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}

func getChildrenCompletions(conn *zk.Conn, root string) []prompt.Suggest {
	if value, ok := suggestCache.get(root); ok {
		return value
	}

	children, _, err := conn.Children(root)
	if err != nil || len(children) == 0 {
		return []prompt.Suggest{}
	}

	s := make([]prompt.Suggest, len(children))
	for i, child := range children {
		p := "/"
		if root == "/" {
			p = fmt.Sprintf("/%s", child)
		} else {
			p = fmt.Sprintf("%s/%s", root, child)
		}
		s[i] = prompt.Suggest{
			Text: p,
		}
	}
	suggestCache.set(root, s)
	return s
}

func splitPath(p string) (root, child string) {
	root, child = path.Split(p)
	root = strings.TrimRight(root, "/")
	if root == "" {
		root = "/"
	}
	return
}
