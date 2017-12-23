package zkcli

import (
	"fmt"
	"path"
	"strings"

	"github.com/c-bata/go-prompt"
)

var commands = []prompt.Suggest{
	{Text: "get", Description: "Get data"},
	{Text: "ls", Description: "List children"},
	{Text: "create", Description: "Create a node"},
	{Text: "set", Description: "Update a node"},
	{Text: "delete", Description: "Delete a node"},
	{Text: "close", Description: "Close connection"},
	{Text: "connect", Description: "Connect servers"},
	{Text: "addauth", Description: "Add auth info"},
	{Text: "exit", Description: "Exit this program"},
}

var suggestCache = newSuggestCache()

func GetCompleter(cmd *Cmd) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		if d.TextBeforeCursor() == "" {
			return []prompt.Suggest{}
		}
		args := strings.Split(d.TextBeforeCursor(), " ")
		return argumentsCompleter(excludeOptions(args), cmd)
	}
}

func argumentsCompleter(args []string, cmd *Cmd) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}

	first := args[0]
	switch first {
	case "get", "ls", "create", "set", "delete":
		p := args[1]
		if len(args) > 2 {
			switch first {
			case "create", "set":
				if len(args) < 4 {
					return []prompt.Suggest{
						{Text: "data"},
					}
				}
			default:
				return []prompt.Suggest{}
			}
		}
		root, _ := splitPath(p)
		return prompt.FilterContains(getChildrenCompletions(cmd, root), p, true)
	case "connect":
		servers := args[1]
		if servers == "" {
			return []prompt.Suggest{
				{Text: "host:port"},
			}
		}
	case "addauth":
		scheme := args[1]
		if len(args) > 2 {
			if len(args) == 3 {
				return []prompt.Suggest{
					{Text: "auth"},
				}
			}
			return []prompt.Suggest{}
		}
		return prompt.FilterContains([]prompt.Suggest{
			{Text: "digest"},
		}, scheme, true)

	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}

func getChildrenCompletions(cmd *Cmd, root string) []prompt.Suggest {
	if value, ok := suggestCache.get(root); ok {
		return value
	}

	if !cmd.connected() {
		return []prompt.Suggest{}
	}

	children, _, err := cmd.Conn.Children(root)
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
